package handler

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"

	"dreampicai/pkg/db"
	"dreampicai/utils"
	"dreampicai/view/auth"
)

func AuthCallback(w http.ResponseWriter, r *http.Request) error {
	// that if potetnial forever loop, i would handle it somehow, like frontend side validation
	// with redirect to some error page with instructions for users, but since my only user is
	// me, i'm ok with that
	if len(r.URL.Query()) == 0 {
		slog.Info("[AuthCallback]", "err", "no query params")
		return auth.CallbackScript().Render(r.Context(), w)
	}

	accessToken, refreshToken, err := utils.GetTokensFromQuery(r.URL.Query())
	if err != nil {
		slog.Info("[AuthCallback] parsing tokens from query", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return err
	}

	userAuth, err := utils.ParseSupabaseToken(accessToken)
	if err != nil {
		slog.Info("[AuthCallback] parsing access token", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return err
	}

	uuidBytes, err := utils.ToUUIDBytes(userAuth.ID)
	if err != nil {
		slog.Info("[AuthCallback] converting user id to byte[16]", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return err
	}

	account, err := db.Q.AccountGetByUserId(
		r.Context(),
		pgtype.UUID{Bytes: uuidBytes, Valid: true},
	)
	if errors.Is(err, sql.ErrNoRows) {
		account, err = db.Q.AccountCreate(
			r.Context(),
			db.AccountCreateParams{
				UserID:   pgtype.UUID{Bytes: uuidBytes, Valid: true},
				Username: userAuth.Email,
			},
		)
	}
	if err != nil {
		slog.Info("[AuthCallback] getting/creating account", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return err
	}

	slog.Info(
		"[AuthCallback] success",
		"accountId",
		strconv.Itoa(int(account.ID)),
		"accessToken",
		accessToken,
		"refreshToken",
		refreshToken,
	)
	utils.AddUserAuthCookies(w, accessToken, refreshToken, strconv.Itoa(int(account.ID)))
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}
