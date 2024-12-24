package handler

import (
	"net/http"

	"dreampicai/view"
)

func GenerateView(w http.ResponseWriter, r *http.Request) error {
	return view.Generate().Render(r.Context(), w)
}

func ListImages(w http.ResponseWriter, r *http.Request) error {
	return nil
}
