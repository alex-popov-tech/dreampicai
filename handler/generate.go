package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"dreampicai/domain"
	"dreampicai/pkg/db"
	"dreampicai/pkg/replicate"
	"dreampicai/utils"
	"dreampicai/view"
)

func GenerateView(w http.ResponseWriter, r *http.Request) error {
	account := utils.GetAccountFromRequest(r)
	dbImages, err := db.Client.ImageList(r.Context(), pgtype.Int4{Int32: account.ID, Valid: true})
	if err != nil {
		return err
	}

	images := make([]domain.Image, len(dbImages))
	for i, dbImage := range dbImages {
		images[i] = domain.Image{
			ID:     dbImage.ID,
			Status: dbImage.Status,
			Prompt: dbImage.Prompt,
		}
		if dbImage.Url.Valid {
			images[i].Url = dbImage.Url.String
		}
	}

	data := view.GenerateViewData{
		Images: images,
		GenerateFormData: view.GenerateFormData{
			SelectedModel:  toSelectOption(domain.DEFAULT_MODEL),
			Models:         toSelectOptions(domain.MODELS),
			Prompt:         domain.DEFAULT_PROMPT,
			NegativePrompt: domain.DEFAULT_NEGATIVE_PROMPT,
			Count:          domain.DEFAULT_COUNT,
		},
	}
	return view.Generate(data).Render(r.Context(), w)
}

func Generate(w http.ResponseWriter, r *http.Request) error {
	var model domain.ReplicateModel = r.FormValue("model")
	prompt := r.FormValue("prompt")
	negativePrompt := r.FormValue("negative_prompt")
	countStr := r.FormValue("count")
	if _, err := strconv.Atoi(countStr); prompt == "" || negativePrompt == "" || countStr == "" ||
		err != nil {
		time.Sleep(5 * time.Second)
		return fmt.Errorf(
			"Missing/invalid one of required body params, prompt:%s, negative_prompt:%s, count:%s",
			prompt,
			negativePrompt,
			countStr,
		)
	}
	count, _ := strconv.Atoi(countStr)

	for i := 0; i < count; i++ {
		prediction, err := replicate.Client.CreatePrediction(
			r.Context(),
			model,
			replicate.PredictionInput{
				"prompt":              prompt,
				"negative_prompt":     negativePrompt,
				"num_outputs":         count,
				"width":               960,
				"height":              1280,
				"output_quality":      100,
				"prompt_strength":     0.8,
				"num_inference_steps": 10,
			},
			&replicate.Webhook{
				URL:    os.Getenv("REPLICATE_WEBHOOK"),
				Events: []replicate.WebhookEventType{"completed"},
			},
			false,
		)
		if err != nil {
			slog.Info("[Generate] creating prediction", "err", err)
			return err
		}

		_, err = db.Client.ImageCreate(r.Context(), db.ImageCreateParams{
			ProviderID:     prediction.ID,
			OwnerID:        pgtype.Int4{Int32: utils.GetAccountFromRequest(r).ID, Valid: true},
			Status:         "started",
			Prompt:         prompt,
			Model:          model,
			NegativePrompt: negativePrompt,
		})
		if err != nil {
			slog.Info("[Generate] inserting image", "err", err)
			return err
		}
	}

	slog.Info("[Generate] success")
	w.Header().Add("HX-Trigger", "refresh")
	return view.GenerateForm(view.GenerateFormData{
		Prompt:         prompt,
		NegativePrompt: negativePrompt,
		Count:          countStr,
		SelectedModel:  toSelectOption(model),
		Models:         toSelectOptions(domain.MODELS),
	}).
		Render(r.Context(), w)
}

func GeneratedWebhook(w http.ResponseWriter, r *http.Request) error {
	_, err := replicate.ValidateWebhookRequest(r, os.Getenv("REPLICATE_SECRET"))
	if err != nil {
		slog.Info("[GenerateWebhook] Validations webhook", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Webhook validation error",
		})
		return err
	}

	parsed := map[string]interface{}{}
	err = json.NewDecoder(r.Body).Decode(&parsed)
	defer r.Body.Close()
	if err != nil {
		slog.Info("[GenerateWebhook] Parsing json", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Parsing json error",
		})
		return err
	}

	id, ok := parsed["id"].(string)
	if !ok {
		slog.Info("[GenerateWebhook] Missing prediction id")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing prediction id",
		})
		return err
	}
	rawStatus, ok := parsed["status"].(string)
	if !ok {
		slog.Info("[GenerateWebhook] Missing prediction status")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing prediction status",
		})
		return err
	}
	status, err := getImageStatusFromString(rawStatus)
	if err != nil {
		slog.Info("[GenerateWebhook] Wrong prediction status", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Wrong prediction status",
		})
		return err
	}
	url, err := getImageUrlFromReplicate(parsed)
	if err != nil {
		slog.Info("[GenerateWebhook] Missing output image url", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing output image url",
		})
		return err
	}

	db.Client.ImageUpdate(r.Context(), db.ImageUpdateParams{
		ProviderID: id,
		Status:     status,
		Url:        pgtype.Text{String: url, Valid: true},
	})
	return err
}

func getImageStatusFromString(status string) (db.ImageStatus, error) {
	switch status {
	case "succeeded":
		return db.ImageStatusSucceeded, nil
	case "failed":
		return db.ImageStatusFailed, nil
	default:
		return "", fmt.Errorf("invalid image status: %s", status)
	}
}

func getImageUrlFromReplicate(responseBody map[string]interface{}) (string, error) {
	if outputArr, ok := responseBody["output"].([]interface{}); ok && len(outputArr) > 0 &&
		outputArr[0] != nil {
		if url, ok := outputArr[0].(string); ok && url != "" {
			return url, nil
		}
	}
	return "", fmt.Errorf("failed to get image url from replicate")
}

func toSelectOptions(models []domain.ReplicateModel) []view.SelectOption {
	options := make([]view.SelectOption, len(models))
	for i, model := range models {
		options[i] = toSelectOption(model)
	}
	return options
}

func toSelectOption(model domain.ReplicateModel) view.SelectOption {
	switch model {
	case domain.REPLICATE_MODEL_PLAYGROUND:
		return view.SelectOption{
			Value: model,
			Text:  "playgroundai/playground-v2.5-1024px-aesthetic",
		}
	case domain.REPLICATE_MODEL_KADNINSKY:
		return view.SelectOption{
			Value: model,
			Text:  "ai-forged/kadninsky-2.0",
		}
	case domain.REPLICATE_MODEL_PROTEUS:
		return view.SelectOption{
			Value: model,
			Text:  "datacte / proteus-v0.3",
		}
	default:
		return view.SelectOption{
			Value: model,
			Text:  "Unknown Model",
		}
	}
}
