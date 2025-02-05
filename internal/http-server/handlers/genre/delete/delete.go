package delete

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	resp "github.com/stepan41k/testMidlware/internal/lib/api/response"
	"github.com/stepan41k/testMidlware/internal/lib/logger/sl"
	"github.com/stepan41k/testMidlware/internal/storage"
)

type GenreDeleter interface {
	DeleteGenre(genre string) error
}

func New(log *slog.Logger, genreDeleter GenreDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.genre.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		genre := chi.URLParam(r, "genre")
		if genre == "" {
			log.Info("genre is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err := genreDeleter.DeleteGenre(genre)
		if errors.Is(err, storage.ErrGenreNotFound) {
			log.Info("genre not found")

			render.JSON(w, r, resp.Error("genre not found"))

			return
		}

		if err != nil {
			log.Error("failed to delete genre", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("genre deleted", slog.String("genre", genre))

		render.JSON(w, r, resp.Response{
			Status: resp.StatusOK,
		})
	}
}