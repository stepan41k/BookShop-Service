package delete

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/stepan41k/testMidlware/internal/lib/api/response"
	"github.com/stepan41k/testMidlware/internal/lib/logger/sl"
	"github.com/stepan41k/testMidlware/internal/storage"
)

type Request struct {
	Name string `json:"name" validate:"required"`
}

type Response struct {
	resp.Response
	Name string `json:"name"`
}

type BookDelter interface {
	DeleteBook(name string) error
}

func New(log *slog.Logger, bookDeleter BookDelter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.book.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		name := chi.URLParam(r, "name")
		if name == "" {
			log.Info("name is empty")

			render.JSON(w, r, resp.Error("name is empty"))

			return
		}

		err := bookDeleter.DeleteBook(name)
		if errors.Is(err, storage.ErrBookNotFound) {
			log.Error("book not found", "name", name)

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		if err != nil {
			log.Error("failed to delete book", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("delete book", slog.String("name", name))

		render.JSON(w, r, resp.Response{Status: resp.StatusOK})
	}
}