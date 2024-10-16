package utils

import (
	"context"
	"dreampicai/model"
	"dreampicai/pkg/supabase"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/golang-jwt/jwt/v5"
)

func GetTokensFromQuery(values url.Values) (accessToken string, refreshToken string, err error) {
	accessToken = values.Get("access_token")
	refreshToken = values.Get("refresh_token")
	if accessToken == "" || refreshToken == "" {
		return "", "", fmt.Errorf("Cannot get 'access_token' or 'refresh_token' from query: %v", values)
	}
	return accessToken, refreshToken, nil
}

func MakeRoute(handler func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic occurred in handler %s %s\nError: %v", r.Method, r.URL.Path, err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(fmt.Sprintf("Sorry, something went VERY wrong\n%v", err)))
			}
		}()

		err := handler(w, r)

		if err != nil {
			log.Printf("Handler unhanded error %s %s\nError: %v", r.Method, r.URL.Path, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Sorry, something went wrong\n%v", err.Error())))
		}
	})
}

func UserAuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atCookie, err := r.Cookie("at")
		at := atCookie.Value
		rtCookie, err := r.Cookie("rt")
		rt := rtCookie.Value
		requestWithEmptyUserContext := r.WithContext(context.WithValue(r.Context(), model.UserContextKey, model.User{}))
		if err != nil {
			log.Printf("UserAuthMiddleware, cookie not found: %v\n", err)
			handler.ServeHTTP(w, requestWithEmptyUserContext)
			return
		}
		if at == "" {
			log.Printf("UserAuthMiddleware, access_token cookie is empty: %v\n", err)
			handler.ServeHTTP(w, requestWithEmptyUserContext)
			return
		}
		if rt == "" {
			log.Printf("UserAuthMiddleware, refresh_token cookie is empty: %v\n", err)
			handler.ServeHTTP(w, requestWithEmptyUserContext)
			return
		}
		user, err := ParseSupabaseToken(at)
		if errors.Is(err, jwt.ErrTokenExpired) {
			authDetails, err := supabase.Client.Auth.RefreshUser(context.Background(), at, rt)
			if err != nil {
				log.Printf("UserAuthMiddleware, access_token is expired and cannot be refreshed: %v\n", err)
				handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), model.UserContextKey, user)))
				return
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
			handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), model.UserContextKey, model.User{Email: authDetails.User.Email, ID: authDetails.User.ID, IsLoggedIn: true})))
			return
		}
		if err != nil {
			log.Printf("UserAuthMiddleware, jwt parsing error: %v\n", err)
			handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), model.UserContextKey, user)))
			return
		}
		handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), model.UserContextKey, user)))
	})
}

func AuthProtectedMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user, ok := r.Context().Value(model.UserContextKey).(model.User); ok && user.IsLoggedIn {
			handler.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/signin?redirect=%s", r.URL.Path), http.StatusFound)
	})
}
