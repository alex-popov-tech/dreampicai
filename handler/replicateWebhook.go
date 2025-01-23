package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"dreampicai/pkg/db"
	"dreampicai/pkg/replicate"
)

func GeneratedWebhook(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
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
	status, err := toStatus(rawStatus)
	if err != nil {
		slog.Info("[GenerateWebhook] Wrong prediction status", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Wrong prediction status",
		})
		return err
	}
	if status == db.ImageStatusFailed {
		_, err = db.Q.ImageUpdate(context.Background(), db.ImageUpdateParams{
			ProviderID: id,
			Status:     status,
		})
		return nil
	}

	urls, err := getImageUrlsFromReplicate(parsed)
	if err != nil {
		slog.Info("[GenerateWebhook] Missing output image url", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Missing output image url",
		})
		return err
	}

	for i, url := range urls {
		ext := path.Ext(url)
		filename := fmt.Sprintf("%s_%d%s", id, i, string(ext))
		go func() {
			err = writeImage(url, path.Join(os.Getenv("IMAGES_DIR"), filename))
			if err != nil {
				slog.Info("[Generate] writing cache", "err", err)
			}
			_, err = db.Q.ImageUpdate(context.Background(), db.ImageUpdateParams{
				ProviderID: id,
				Status:     status,
				Filename:   pgtype.Text{String: filename, Valid: true},
			})
			if err != nil {
				slog.Info("[Generate] saving to db", "err", err)
			}
			slog.Info("[GenerateWebhook] Saved file to db and cache", "filename", filename)
		}()
	}

	return nil
}

func toStatus(status string) (db.ImageStatus, error) {
	switch status {
	case "succeeded":
		return db.ImageStatusSucceeded, nil
	case "failed":
		return db.ImageStatusFailed, nil
	default:
		return "", fmt.Errorf("invalid image status: %s", status)
	}
}

func getImageUrlsFromReplicate(responseBody map[string]interface{}) ([]string, error) {
	if outputArr, ok := responseBody["output"].([]interface{}); ok && len(outputArr) > 0 {
		result := make([]string, len(outputArr))
		for i, output := range outputArr {
			result[i] = output.(string)
		}
		return result, nil
	}
	return []string{}, fmt.Errorf("failed to get image urls from replicate")
}

func writeImage(url, filepath string) error {
	maxRetries := 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		if err := downloadImage(url, filepath); err != nil {
			lastErr = err
			continue
		}
		return nil
	}
	return fmt.Errorf("failed after %d attempts: %v", maxRetries, lastErr)
}

func downloadImage(url, filepath string) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	return err
}
