package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"dreampicai/pkg/supabase"
	"dreampicai/utils"
	"dreampicai/view/auth"
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
	fmt.Sprintf("hello mom! %s", "hello mom")

	loginData := auth.SigninData{Email: email, Password: password}
	loginErrors := auth.SigninErrors{Email: emailErrors, Password: passwordErrors}

	if len(emailErrors) > 0 || len(passwordErrors) > 0 {
		return auth.SigninForm(loginData, loginErrors, redirect).Render(r.Context(), w)
	}

	authDetails, err := supabase.Client.Auth.SignIn(
		context.Background(),
		supabase.UserCredentials{Email: email, Password: password},
	)
	if err != nil {
		loginErrors.Others = []error{err}
		return auth.SigninForm(loginData, loginErrors, redirect).Render(r.Context(), w)
	}

	utils.AddAuthCookies(w, authDetails.AccessToken, authDetails.RefreshToken)
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
