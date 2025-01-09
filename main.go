package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"dreampicai/handler"
	"dreampicai/pkg/db"
	"dreampicai/pkg/replicate"
	"dreampicai/pkg/supabase"
	"dreampicai/utils"
)

//go:embed public
var content embed.FS

func initialize() error {
	env, err := utils.ValidateEnv()
	if err != nil {
		return fmt.Errorf("Error loading environment variables \n%v\n", err)
	}

	_ = supabase.InitClient(env.SupabaseProjectURL, env.SupabaseServiceSecretKey)

	_, err = db.InitClient(env.DatabaseDirectURL)
	if err != nil {
		return fmt.Errorf("Error creating database client \n%v\n", err)
	}

	_, err = replicate.InitClient(env.ReplicateToken)
	if err != nil {
		return fmt.Errorf("Error creating replicate client \n%v\n", err)
	}

	return nil
}

func main() {
	if err := initialize(); err != nil {
		log.Fatal(err)
	}
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)

	mux.Handle("/*", http.StripPrefix("/public/", http.FileServerFS(os.DirFS("public"))))

	mux.Handle("GET /signin", utils.MakeRoute(handler.SigninView))
	mux.Handle("POST /signin", utils.MakeRoute(handler.Signin))
	mux.Handle("/signin/github", utils.MakeRoute(handler.SigninWithGithub))

	mux.Handle("GET /signup", utils.MakeRoute(handler.SignupView))
	mux.Handle("POST /signup", utils.MakeRoute(handler.Signup))

	mux.Handle("DELETE /signout", utils.MakeRoute(handler.Signout))
	mux.Handle("/auth/callback", utils.MakeRoute(handler.AuthCallback))

	mux.Group(func(mux chi.Router) {
		mux.Use(utils.WithUser)
		mux.Handle("GET /", utils.MakeRoute(handler.HomeView))

		mux.Group(func(mux chi.Router) {
			mux.Use(utils.Protected)
			mux.Handle("GET /settings", utils.MakeRoute(handler.SettingsView))
		})
		mux.Handle("GET /generate", utils.MakeRoute(handler.GenerateView))
		mux.Handle("POST /generate", utils.MakeRoute(handler.Generate))
		mux.Handle("GET /images/{id}", utils.MakeRoute(handler.GetImage))
	})

	mux.Handle("POST /webhook/replicate", utils.MakeRoute(handler.GeneratedWebhook))

	fmt.Println("Listening on port", os.Getenv("PORT"))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), mux); err != nil {
		log.Fatal("oops", err)
	}
}
