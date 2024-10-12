package utils

import (
	"context"
	"dreampicai/model"
	"fmt"
	"log"
	"net/http"
)

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
		user := model.User{}
		if err != nil {
			log.Printf("UserAuthMiddleware, cookie not found: %v\n", err)
			handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), model.UserContextKey, user)))
			return
		}
		if atCookie.Value == "" {
			log.Printf("UserAuthMiddleware, cookie is empty: %v\n", err)
			handler.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), model.UserContextKey, user)))
			return
		}
		user, err = ParseSupabaseToken(atCookie.Value)
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
		http.Redirect(w, r, "/signin", http.StatusFound)
	})
}
