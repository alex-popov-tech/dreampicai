package utils

import (
	"context"

	"dreampicai/domain"
)


func GetAccount(ctx context.Context) domain.Account {
	auth, ok := ctx.Value(domain.AccountContextKey).(domain.Account)
	if !ok {
		return domain.Account{}
	}
	return auth
}
