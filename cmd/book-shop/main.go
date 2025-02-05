package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stepan41k/testMidlware/internal/config"
	"github.com/stepan41k/testMidlware/internal/lib/logger/sl"
	"github.com/stepan41k/testMidlware/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev = "dev"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("info messages are enabled")
	log.Debug("debug messages are enabled")
	log.Error("error messages are enabled")

	storage, err := postgres.New(fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
	}
	
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use() //TODO: Logger
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/book", func(r chi.Router) {
		r.Post("/", saveBook.New(log, storage))
		r.Delete("/{book}", deleteBook.New(log, storage))
		r.Get("/{book}", takeBook.New(log, storage))
	})

	router.Route("/genre", func(r chi.Router) {
		r.Post("/", saveGenre.New(log, storage))
		r.Delete("/{book}", deleteGenre.New(log, storage))
	})

	router.Route("/author", func(r chi.Router) {
		r.Post("/", saveAuthor.New(log, storage))
		r.Delete("/{book}", deleteAuthor.New(log, storage))
	})


	log.Info("starting server", slog.String("adress", cfg.Adress))

	srv := http.Server{
		Addr: cfg.Adress,
		Handler: router,
		ReadTimeout: cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout: cfg.HttpServer.Idle_timeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

}

func setupLogger(env string) *slog.Logger{
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}