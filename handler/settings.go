package handler

import (
	"dreampicai/view/settings"
	"net/http"
)

func SettingsView(w http.ResponseWriter, r *http.Request) error {
	return settings.Settings().Render(r.Context(), w)
}
