package handler

import (
	"context"
	"dreampicai/pkg/supabase"
	"dreampicai/utils"
	"dreampicai/view/auth"
	"net/http"
)

func SigninView(w http.ResponseWriter, r *http.Request) error {
	return auth.Signin().Render(r.Context(), w)
}

func Signin(w http.ResponseWriter, r *http.Request) error {
	email := r.FormValue("email")
	password := r.FormValue("password")
	emailErrors := utils.ValidateEmail(email)
	passwordErrors := utils.ValidatePassword(password)

	loginData := auth.SigninData{Email: email, Password: password}
	loginErrors := auth.SigninErrors{Email: emailErrors, Password: passwordErrors}

	if len(emailErrors) > 0 || len(passwordErrors) > 0 {
		return auth.SigninForm(loginData, loginErrors).Render(r.Context(), w)
	}

	authDetails, err := supabase.Client.Auth.SignIn(context.Background(), supabase.UserCredentials{Email: email, Password: password})

	if err != nil {
		loginErrors.Others = []error{err}
		return auth.SigninForm(loginData, loginErrors).Render(r.Context(), w)
	}

	atCookie := http.Cookie{
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Value:    authDetails.AccessToken,
	}
	rtCookie := http.Cookie{
		Name:     "rt",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Value:    authDetails.RefreshToken,
	}
	http.SetCookie(w, &atCookie)
	http.SetCookie(w, &rtCookie)

	w.Header().Add("HX-Redirect", "/")

	return nil
}
