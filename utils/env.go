package utils

import (
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Env struct {
	Port                     uint   `env:"PORT" validate:"required,min=1024,max=65535"`
	DatabaseURL              string `env:"DATABASE_URL" validate:"required,url"`
	SupabaseProjectURL       string `env:"SUPABASE_PROJECT_URL" validate:"required,url"`
	SupabasePublicKey        string `env:"SUPABASE_PUBLIC_KEY" validate:"required"`
	SupabaseServiceSecretKey string `env:"SUPABASE_SERVICE_SECRET_KEY" validate:"required"`
	SupabaseJWTSecret        string `env:"SUPABASE_JWT_SECRET" validate:"required"`
}

func ValidateEnv() (Env, error) {
	err := godotenv.Load()
	if err != nil {
		return Env{}, fmt.Errorf("Error loading .env file\n%s", err.Error())
	}

	e := Env{}
	if err := env.Parse(&e); err != nil {
		return Env{}, fmt.Errorf("Error parsing environment variables:", err)
	}

	validate := validator.New()
	err = validate.Struct(e)
	if err != nil {
		return Env{}, fmt.Errorf("Error validating environment variables:", err)
	}

	return e, nil
}
