package handler

import (
	"context"
	"dreampicai/pkg/supabase"
	"dreampicai/utils"
	"dreampicai/view/auth"
	"fmt"
	"net/http"
	"os"
)

func SigninView(w http.ResponseWriter, r *http.Request) error {
	return auth.Signin(r.URL.Query().Get("redirect")).Render(r.Context(), w)
}

func Signin(w http.ResponseWriter, r *http.Request) error {
	email := r.FormValue("email")
	password := r.FormValue("password")
	redirect := r.FormValue("redirect")
	emailErrors := utils.ValidateEmail(email)
	passwordErrors := utils.ValidatePassword(password)

	loginData := auth.SigninData{Email: email, Password: password}
	loginErrors := auth.SigninErrors{Email: emailErrors, Password: passwordErrors}

	if len(emailErrors) > 0 || len(passwordErrors) > 0 {
		return auth.SigninForm(loginData, loginErrors, redirect).Render(r.Context(), w)
	}

	authDetails, err := supabase.Client.Auth.SignIn(context.Background(), supabase.UserCredentials{Email: email, Password: password})

	if err != nil {
		loginErrors.Others = []error{err}
		return auth.SigninForm(loginData, loginErrors, redirect).Render(r.Context(), w)
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

	fmt.Println("signin success, redirect is", redirect)
	if redirect != "" {
		w.Header().Add("HX-Redirect", redirect)
	} else {
		w.Header().Add("HX-Redirect", "/")
	}

	return nil
}

func SigninWithGithub(w http.ResponseWriter, r *http.Request) error {
	details, err := supabase.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "github",
		RedirectTo: os.Getenv("GITHUB_AUTH_REDIRECT"),
	})
	if err != nil {
		return err
	}

	http.Redirect(w, r, details.URL, http.StatusFound)
	return nil
}
