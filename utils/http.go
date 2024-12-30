package utils

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"dreampicai/domain"
)

func GetAccountFromRequest(r *http.Request) domain.Account {
	return r.Context().Value(domain.AccountContextKey).(domain.Account)
}

func AddUserAuthCookies(w http.ResponseWriter, accessToken, refreshToken, accountId string) {
	http.SetCookie(w, &http.Cookie{
		Name:     domain.AccessTokenCookieKey,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Value:    accessToken,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     domain.RefreshTokenCookieKey,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Value:    refreshToken,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     domain.AccountIdCookieKey,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Value:    accountId,
	})
}

func CleanAllCookies(w http.ResponseWriter, r *http.Request) {
	// Get all cookies from the request
	cookies := r.Cookies()

	// For each cookie, create a new cookie with the same name but expired
	for _, cookie := range cookies {
		// Create a new cookie with the same name
		expiredCookie := &http.Cookie{
			Name:     cookie.Name,
			Value:    "",              // Empty the value
			Path:     "/",             // Cover all paths
			Domain:   cookie.Domain,   // Use same domain
			Expires:  time.Unix(0, 0), // Set to epoch time (expired)
			MaxAge:   -1,              // Delete cookie immediately
			Secure:   cookie.Secure,   // Maintain secure flag
			HttpOnly: cookie.HttpOnly, // Maintain HttpOnly flag
			SameSite: cookie.SameSite, // Maintain SameSite policy
		}

		// Set the expired cookie in response
		http.SetCookie(w, expiredCookie)
	}
}

func GetTokensFromQuery(values url.Values) (accessToken string, refreshToken string, err error) {
	accessToken = values.Get("access_token")
	refreshToken = values.Get("refresh_token")
	if accessToken == "" || refreshToken == "" {
		return "", "", fmt.Errorf(
			"Cannot get 'access_token' or 'refresh_token' from query: %v",
			values,
		)
	}
	return accessToken, refreshToken, nil
}

func MakeRoute(handler func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			log.Printf("Handler unhanded error %s %s\nError: %v", r.Method, r.URL.Path, err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Something went wrong: %v", err)
		}
	})
}
