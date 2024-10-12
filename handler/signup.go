package handler

import (
	"dreampicai/pkg/supabase"
	"dreampicai/utils"
	"dreampicai/view/auth"
	"errors"
	"net/http"

	"github.com/supabase-community/gotrue-go/types"
)

func SignupView(w http.ResponseWriter, r *http.Request) error {
	return auth.Signup().Render(r.Context(), w)
}

func Signup(w http.ResponseWriter, r *http.Request) error {
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirmPassword")

	emailErrors := utils.ValidateEmail(email)
	passwordErrors := utils.ValidatePassword(password)
	confirmPasswordErrors := utils.ValidatePassword(confirmPassword)
	if password != confirmPassword {
		confirmPasswordErrors = append(confirmPasswordErrors, errors.New("Passwords do not match"))
	}

	loginData := auth.SignupData{Email: email, Password: password}
	loginErrors := auth.SignupErrors{Email: emailErrors, Password: passwordErrors, ConfirmPassword: confirmPasswordErrors}

	if len(emailErrors) > 0 || len(passwordErrors) > 0 {
		return auth.SignupForm(loginData, loginErrors).Render(r.Context(), w)
	}

	_, err := supabase.Client.Auth.Signup(types.SignupRequest{Email: email, Password: password, Data: map[string]interface{}{"redirectTo": "https://dreampicai.com/"}})
	if err != nil {
		loginErrors.Others = []error{supabase.TryGerSupabaseErrorMessage(err)}
		return auth.SignupForm(loginData, loginErrors).Render(r.Context(), w)
	}

	return auth.SignupSuccessMessage(email).Render(r.Context(), w)
}
