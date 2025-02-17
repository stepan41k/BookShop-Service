package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stepan41k/testMidlware/internal/config"

	ssogrpc "github.com/stepan41k/testMidlware/internal/clients/sso/grpc"

	deleteAuthor "github.com/stepan41k/testMidlware/internal/http-server/handlers/author/delete"
	saveAuthor "github.com/stepan41k/testMidlware/internal/http-server/handlers/author/save"
	deleteBook "github.com/stepan41k/testMidlware/internal/http-server/handlers/book/delete"
	saveBook "github.com/stepan41k/testMidlware/internal/http-server/handlers/book/save"
	deleteGenre "github.com/stepan41k/testMidlware/internal/http-server/handlers/genre/delete"
	saveGenre "github.com/stepan41k/testMidlware/internal/http-server/handlers/genre/save"
	"github.com/stepan41k/testMidlware/internal/lib/logger/handlers/slogpretty"
	"github.com/stepan41k/testMidlware/internal/lib/logger/sl"
	eventsender "github.com/stepan41k/testMidlware/internal/services/event-sender"
	"github.com/stepan41k/testMidlware/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("info messages are enabled")
	log.Debug("debug messages are enabled")
	log.Error("error messages are enabled")

	ssoClient, err := ssogrpc.New(
		context.Background(),
		log,
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {
		log.Error("failed to init sso client", sl.Err(err))
		os.Exit(1)
	}

	ssoClient.IsAdmin(context.Background(), 1)

	storage, err := postgres.New(fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, os.Getenv("DB_PASSWORD"), cfg.SSLMode))
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
		//r.Get("/{book}", takeBook.dNew(log, storage))
		//r.Get("/", takeBooks.New(log, storage))
		//r.Patch("/{book}", updateBook.New(log, storage))
	})

	router.Route("/genre", func(r chi.Router) {
		r.Post("/", saveGenre.New(log, storage))
		r.Delete("/{genre}", deleteGenre.New(log, storage))
	})

	router.Route("/author", func(r chi.Router) {
		r.Post("/", saveAuthor.New(log, storage))
		r.Delete("/{author}", deleteAuthor.New(log, storage))
	})

	log.Info("starting server", slog.String("adress", cfg.Adress))

	srv := http.Server{
		Addr:         cfg.Adress,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.Idle_timeout,
	}

	sender := eventsender.New(storage, log)
	sender.StartProcessEvents(context.Background(), 5*time.Second)

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
