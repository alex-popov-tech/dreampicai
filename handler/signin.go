package handler

import (
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

	session, err := supabase.Client.Auth.SignInWithEmailPassword(email, password)
	if err != nil {
		loginErrors.Others = []error{supabase.TryGerSupabaseErrorMessage(err)}
		return auth.SigninForm(loginData, loginErrors).Render(r.Context(), w)
	}

	atCookie := http.Cookie{
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Value:    session.AccessToken,
	}
	rtCookie := http.Cookie{
		Name:     "rt",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Value:    session.RefreshToken,
	}
	http.SetCookie(w, &atCookie)
	http.SetCookie(w, &rtCookie)

	w.Header().Add("HX-Redirect", "/")

	return nil
}
