package utils

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"

	"dreampicai/domain"
	"dreampicai/pkg/db"
	"dreampicai/pkg/supabase"
)

func WithUser(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestWithEmptyUserContext := r.WithContext(
			context.WithValue(r.Context(), domain.AuthContextKey, domain.Auth{}),
		)

		atCookie, err := r.Cookie("at")
		// if there is no access token cookie, user is logged out
		if err != nil {
			slog.Info("WithUser", "Cannot get 'at' cookie", err)
			handler.ServeHTTP(w, requestWithEmptyUserContext)
			return
		}

		authUser, err := ParseSupabaseToken(atCookie.Value)
		if err == nil {
			slog.Info(
				"WithUser, parsed access token",
				"User id",
				authUser.ID,
				"User email",
				authUser.Email,
			)
			// jwt successfully parsed and active, just put data to context
			handler.ServeHTTP(
				w,
				r.WithContext(
					context.WithValue(
						r.Context(),
						domain.AuthContextKey,
						domain.Auth{ID: authUser.ID, Email: authUser.Email, IsInit: true},
					),
				),
			)
			return
		}

		// if something is wrong with access token, try to refresh it
		rtCookie, err := r.Cookie("rt")
		if err != nil {
			slog.Info("WithUser, parsing refresh token", "error", err)
			handler.ServeHTTP(w, requestWithEmptyUserContext)
			return
		}
		authDetails, err := supabase.Client.Auth.RefreshUser(
			r.Context(),
			atCookie.Value,
			rtCookie.Value,
		)
		if err != nil {
			// all of this is fucked up, just clean cookies and go home
			slog.Info("WithUser, refreshing token", "error", err)
			CleanAllCookies(w, r)
			handler.ServeHTTP(w, requestWithEmptyUserContext)
			return
		}

		slog.Info(
			"WithUser, token refreshed",
			"at",
			authDetails.AccessToken,
			"rt",
			authDetails.RefreshToken,
		)
		AddAuthCookies(w, authDetails.AccessToken, authDetails.RefreshToken)
		handler.ServeHTTP(
			w,
			r.WithContext(
				context.WithValue(
					r.Context(),
					domain.AuthContextKey,
					domain.Auth{
						Email:  authDetails.User.Email,
						ID:     authDetails.User.ID,
						IsInit: true,
					},
				),
			),
		)
	})
}

func UserProtected(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Context().Value(domain.AuthContextKey).(domain.Auth); ok {
			handler.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/signin?redirect=%s", r.URL.Path), http.StatusFound)
	})
}

func WithAccount(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if there is no user, return
		user, ok := r.Context().Value(domain.AuthContextKey).(domain.Auth)
		if !ok {
			handler.ServeHTTP(w, r)
			return
		}
		bytes, err := ToUUIDBytes(user.ID)
		if err != nil {
			log.Printf("WithAccount: user.ID cannot be converted to UUID bytes - %v\n", err)
			handler.ServeHTTP(w, r)
			return
		}

		// if user is there, grab account from db and put into context
		acc, err := db.Client.AccountGetByUserId(
			r.Context(),
			pgtype.UUID{Bytes: bytes, Valid: true},
		)
		if err != nil {
			log.Printf("WithAccount: cannot find account for user.ID '%s' - %v\n", user.ID, err)
			handler.ServeHTTP(w, r)
			return
		}

		handler.ServeHTTP(
			w,
			r.WithContext(context.WithValue(r.Context(), domain.AccountContextKey, acc)),
		)
	})
}
