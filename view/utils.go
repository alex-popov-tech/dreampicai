package view

import (
	"context"
	"dreampicai/domain"
)

func GetUser(ctx context.Context) domain.User {
	user, ok := ctx.Value(domain.UserContextKey).(domain.User)
	if !ok {
		return domain.User{}
	}
	return user
}
