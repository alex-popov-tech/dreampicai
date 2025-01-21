package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"

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

	err = utils.EnsurePathExists(env.ImagesDir)
	if err != nil {
		return fmt.Errorf("Error creating dir for image generation \n%v\n", err)
	}

	_ = supabase.InitClient(env.SupabaseProjectURL, env.SupabaseServiceSecretKey)

	_, err = db.InitClient(env.DatabasePoolURL)
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

	// server generated images
	generatedImagesRoutePath := "/" + os.Getenv("IMAGES_DIR") + "/"
	fileServer := http.FileServer(http.Dir(os.Getenv("IMAGES_DIR")))
	mux.Handle(generatedImagesRoutePath+"*", http.StripPrefix(generatedImagesRoutePath,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age=3600")
			fileServer.ServeHTTP(w, r)
		})))

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
			mux.Handle("GET /generate", utils.MakeRoute(handler.GenerateView))
			mux.Handle("POST /generate", utils.MakeRoute(handler.Generate))
			mux.Handle("GET /images/{id}", utils.MakeRoute(handler.GetImage))
			mux.Handle("GET /images", utils.MakeRoute(handler.GetImages))
		})
	})

	mux.Handle("POST /webhook/replicate", utils.MakeRoute(handler.GeneratedWebhook))

	go startPingingSupabase()

	log.Default().Printf("Listening on port %s\n", os.Getenv("PORT"))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), mux); err != nil {
		log.Fatal("oops", err)
	}
}

func startPingingSupabase() {
	for {
		// random duration between 24 and 57 hours
		duration := time.Duration(24+rand.Intn(33)) * time.Hour
		slog.Info("[PingingSupabase] sleeping for", "hours", int64(duration/time.Hour))
		refreshSupabase()
		time.Sleep(duration)
	}
}

func refreshSupabase() {
	timeout := 30 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	doneCh := make(chan bool)
	accountsCount := -1
	imagesCount := -1

	go func() {
		err := db.C.QueryRow(ctx, "select count(*) from accounts").Scan(&accountsCount)
		if err != nil || accountsCount <= 0 {
			slog.Info("[RefreshSupabase] counting accounts", "err", err)
		}

		err = db.C.QueryRow(ctx, "select count(*) from images").Scan(&imagesCount)
		if err != nil || imagesCount <= 0 {
			slog.Info("[RefreshSupabase] counting images", "err", err)
		}
		doneCh <- true
	}()

	select {
	case <-doneCh:
		slog.Info("[RefreshSupabase]", "imagesCount", imagesCount, "accountsCount", accountsCount)
	case <-ctx.Done():
		slog.Info("[RefreshSupabase] Timed out!")
	}
}
