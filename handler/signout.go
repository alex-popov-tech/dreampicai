package handler

import (
	"net/http"

	"dreampicai/utils"
)

func Signout(w http.ResponseWriter, r *http.Request) error {
	utils.CleanAllCookies(w, r)
	w.Header().Add("HX-Redirect", "/")
	return nil
}
