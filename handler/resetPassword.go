package handler

import (
	"net/http"

	"dreampicai/pkg/supabase"
	"dreampicai/utils"
	"dreampicai/view/auth"
)

func ResetPasswordView(w http.ResponseWriter, r *http.Request) error {
	return auth.ResetPassword().Render(r.Context(), w)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) error {
	email := r.FormValue("email")
	emailErrors := utils.ValidateEmail(email)

	data := auth.ResetPasswordData{Email: email}
	errors := auth.ResetPasswordErrors{Email: emailErrors}

	if len(emailErrors) > 0 {
		return auth.ResetPasswordForm(data, errors).Render(r.Context(), w)
	}

	err := supabase.Client.Auth.ResetPasswordForEmail(r.Context(), email)
	if err != nil {
		errors.Email = []error{err}
		return auth.ResetPasswordForm(data, errors).Render(r.Context(), w)
	}

	auth.ResetPasswordSuccessMessage(email).Render(r.Context(), w)
	return nil
}
