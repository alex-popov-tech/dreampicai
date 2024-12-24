package handler

import (
	"net/http"

	"dreampicai/view"
)

func SettingsView(w http.ResponseWriter, r *http.Request) error {
	return view.Settings().Render(r.Context(), w)
}
