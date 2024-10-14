package handler

import (
	"dreampicai/utils"
	"dreampicai/view/auth"
	"net/http"
)

func AuthCallback(w http.ResponseWriter, r *http.Request) error {
	// that if potetnial forever loop, i would handle it somehow, like frontend side validation
	// with redirect to some error page with instructions for users, but since my only user is
	// me, i'm ok with that
	if len(r.URL.Query()) == 0 {
		return auth.CallbackScript().Render(r.Context(), w)
	}
	accessToken, refreshToken, err := utils.GetTokensFromQuery(r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Value:    accessToken,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "rt",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Value:    refreshToken,
	})

	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}
