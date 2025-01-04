package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

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

	slog.Info("[GenerateView]", "images", images)
	data := view.GenerateData{Images: images}
	return view.Generate(data).Render(r.Context(), w)
}

func Generate(w http.ResponseWriter, r *http.Request) error {
	prompt := "Ukrainian woman in national dress, warm color palette, muted colors, detailed, 8k"
	prediction, err := replicate.Client.CreatePrediction(
		r.Context(),
		os.Getenv("REPLICATE_MODEL"),
		replicate.PredictionInput{"prompt": prompt},
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

	dbImage, err := db.Client.ImageCreate(r.Context(), db.ImageCreateParams{
		ProviderID: prediction.ID,
		OwnerID:    pgtype.Int4{Int32: utils.GetAccountFromRequest(r).ID, Valid: true},
		Status:     "started",
		Prompt:     prompt,
	})
	if err != nil {
		slog.Info("[Generate] inserting image", "err", err)
		return err
	}

	image := domain.Image{
		ID:     dbImage.ID,
		Status: dbImage.Status,
		Prompt: dbImage.Prompt,
		Url:    dbImage.Url.String,
	}

	slog.Info("[Generate] success")
	return view.Card(image).Render(r.Context(), w)
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
