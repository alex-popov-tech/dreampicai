package utils

import (
	"dreampicai/model"
	"fmt"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

var SIGNING_METHOD = jwt.SigningMethodHS256

func ParseSupabaseToken(token string) (model.User, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		// should be populated in main.go by dotenv
		return []byte(os.Getenv("SUPABASE_JWT_SECRET")), nil
	})
	if err != nil {
		return model.User{}, err
	}
	user := model.User{
		ID:         (claims["sub"].(string)),
		Email:      (claims["email"].(string)),
		IsLoggedIn: true,
	}
	return user, err
}
