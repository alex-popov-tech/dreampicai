package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
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
	dbImages, err := db.Q.ImageList(r.Context(), pgtype.Int4{Int32: account.ID, Valid: true})
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
		if dbImage.Filename.Valid {
			images[i].Url = path.Join(os.Getenv("IMAGES_DIR"), dbImage.Filename.String)
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
				"num_outputs":         1,
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

		_, err = db.Q.ImageCreate(r.Context(), db.ImageCreateParams{
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
		ToastMessage:   "Prompt Submitted! Be aware that generations could take up to 1-3 min depending on model",
		ToastStatus:    "success",
	}).
		Render(r.Context(), w)
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
