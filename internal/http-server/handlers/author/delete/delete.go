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

type AuthorDeleter interface {
	DeleteAuthor(author string) error
}

func New(log *slog.Logger, authorDeleter AuthorDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.author.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		author := chi.URLParam(r, "author")
		if author == "" {
			log.Info("author is empty")

			render.JSON(w, r, resp.Error("author is empty"))

			return
		}

		err := authorDeleter.DeleteAuthor(author)
		if errors.Is(err, storage.ErrAuthorNotFound) {
			log.Info("author not found", "author", author)

			render.JSON(w, r, resp.Error("author not found"))

			return
		}
		if err != nil {
			log.Error("failed to delete author", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to delete author"))

			return
		}

		log.Info("author deleted", slog.String("author", author))

		render.JSON(w, r, resp.Response{
			Status: resp.StatusOK,
		})
	}
}