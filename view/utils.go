package view

import (
	"context"
	"dreampicai/model"
)

func GetUser(ctx context.Context) model.User {
	user, ok := ctx.Value(model.UserContextKey).(model.User)
	if !ok {
		return model.User{}
	}
	return user
}
