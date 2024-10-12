package handler

import (
	"net/http"
)

func Signout(w http.ResponseWriter, r *http.Request) error {
	atCookie := http.Cookie{Name: "at", Value: ""}
	rtCookie := http.Cookie{Name: "rt", Value: ""}
	http.SetCookie(w, &atCookie)
	http.SetCookie(w, &rtCookie)

	w.Header().Add("HX-Redirect", "/")

	return nil
}
