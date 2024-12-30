package utils

import (
	"errors"
	"fmt"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"

	"dreampicai/pkg/supabase"
)

var SIGNING_METHOD = jwt.SigningMethodHS256

func ParseSupabaseToken(token string) (*supabase.SupabaseAuth, error) {
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
		return nil, err
	}

	appMetadata, ok := claims["app_metadata"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Can't find 'app_metadata' in supabase token")
	}

	return &supabase.SupabaseAuth{
		ID:       (claims["sub"].(string)),
		Email:    (claims["email"].(string)),
		Provider: (appMetadata["provider"].(string)),
	}, nil
}
