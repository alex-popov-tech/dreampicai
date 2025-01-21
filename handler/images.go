package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"dreampicai/domain"
	"dreampicai/pkg/db"
	"dreampicai/utils"
	"dreampicai/view"
)

func GetImages(w http.ResponseWriter, r *http.Request) error {
	account := utils.GetAccountFromRequest(r)
	dbImages, err := db.Q.ImageList(r.Context(), pgtype.Int4{Int32: account.ID, Valid: true})
	if err != nil {
		slog.Info("[GetImages] getting image from db", "err", err)
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Image with such owner not found with",
		})
		return nil
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

	return view.ImageCards(images).Render(r.Context(), w)
}

func GetImage(w http.ResponseWriter, r *http.Request) error {
	account := utils.GetAccountFromRequest(r)
	idAsString := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idAsString, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("invalid image ID: %s", idAsString),
		})
		return nil
	}

	dbImage, err := db.Q.ImageGet(
		r.Context(),
		db.ImageGetParams{ID: int32(id), OwnerID: pgtype.Int4{Int32: account.ID, Valid: true}},
	)
	if err != nil {
		slog.Info("[GetImage] getting image from db", "err", err)
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Image with such owner not found with",
		})
		return nil
	}

	image := domain.Image{
		ID:     dbImage.ID,
		Status: dbImage.Status,
		Prompt: dbImage.Prompt,
	}
	if dbImage.Filename.Valid {
		image.Url = path.Join(os.Getenv("IMAGES_DIR"), dbImage.Filename.String)
	}

	return view.ImageCard(image).Render(r.Context(), w)
}
