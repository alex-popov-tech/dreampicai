package handler

import (
	"net/http"

	"dreampicai/view"
)

func HomeView(w http.ResponseWriter, r *http.Request) error {
	return view.Home().Render(r.Context(), w)
}
