package handler

import (
	"dreampicai/view/home"
	"net/http"
)

func HomeView(w http.ResponseWriter, r *http.Request) error {
	return home.Home().Render(r.Context(), w)
}
