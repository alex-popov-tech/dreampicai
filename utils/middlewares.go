package utils

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"dreampicai/domain"
	"dreampicai/pkg/db"
	"dreampicai/pkg/supabase"
)

func serverRequestWithNoAccount(w http.ResponseWriter, r *http.Request, handler http.Handler) {
	requestWithEmptyUserContext := r.WithContext(
		context.WithValue(r.Context(), domain.AccountContextKey, domain.Account{}),
	)
	handler.ServeHTTP(w, requestWithEmptyUserContext)
}

func getUserAuthFromTokens(accessToken, refreshToken string) (*domain.UserAuth, error) {
	userAuth := &domain.UserAuth{AccessToken: accessToken, RefreshToken: refreshToken}
	supabaseAuth, err := ParseSupabaseToken(accessToken)
	if err == nil {
		userAuth.ID = supabaseAuth.ID
		userAuth.Email = supabaseAuth.Email
		userAuth.Provider = supabaseAuth.Provider
		return userAuth, nil
	}
	if err.Error() != "token is expired" {
		return nil, err
	}

	authDetails, err := supabase.Client.Auth.RefreshUser(
		context.Background(),
		accessToken,
		refreshToken,
	)
	if err != nil {
		return nil, err
	}

	supabaseAuth, err = ParseSupabaseToken(authDetails.AccessToken)
	if err != nil {
		return nil, err
	}

	return &domain.UserAuth{
		ID:           supabaseAuth.ID,
		Email:        supabaseAuth.Email,
		Provider:     supabaseAuth.Provider,
		AccessToken:  authDetails.AccessToken,
		RefreshToken: authDetails.RefreshToken,
	}, nil
}

func WithUser(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atCookie, _ := r.Cookie(domain.AccessTokenCookieKey)
		rtCookie, _ := r.Cookie(domain.RefreshTokenCookieKey)
		accountIdCookie, _ := r.Cookie(domain.AccountIdCookieKey)
		if atCookie == nil || rtCookie == nil || accountIdCookie == nil {
			slog.Info(
				"[WithUser] missing required cookie",
				"at",
				atCookie,
				"rt",
				rtCookie,
				"accountId",
				accountIdCookie,
			)
			serverRequestWithNoAccount(w, r, handler)
			return
		}

		// if there is no access token cookie, user is logged out
		userAuth, err := getUserAuthFromTokens(atCookie.Value, rtCookie.Value)
		if err != nil {
			slog.Info("[WithUser] parsing/refreshing tokens", "err", err)
			CleanAllCookies(w, r)
			serverRequestWithNoAccount(w, r, handler)
			return
		}

		accountId, err := strconv.Atoi(accountIdCookie.Value)
		if err != nil {
			slog.Info("[WithUser] parsing account id", "err", err)
			serverRequestWithNoAccount(w, r, handler)
			return
		}

		account, err := db.Client.AccountGet(r.Context(), int32(accountId))
		if err != nil {
			slog.Info("[WithUser] getting account from db", "err", err)
			serverRequestWithNoAccount(w, r, handler)
			return
		}

		slog.Info("[WithUser] user authenticated", "id", userAuth.ID, "email", userAuth.Email)

		AddUserAuthCookies(w, userAuth.AccessToken, userAuth.RefreshToken, accountIdCookie.Value)

		handler.ServeHTTP(
			w,
			r.WithContext(
				context.WithValue(
					r.Context(),
					domain.AccountContextKey,
					domain.Account{
						ID:       account.ID,
						Email:    userAuth.Email,
						UserAuth: *userAuth,
					},
				),
			),
		)
	})
}

func Protected(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := r.Context().Value(domain.AccountContextKey).(domain.Account); ok {
			handler.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/signin?redirect=%s", r.URL.Path), http.StatusFound)
	})
}
