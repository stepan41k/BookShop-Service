package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/stepan41k/testMidlware/internal/domain"
	resp "github.com/stepan41k/testMidlware/internal/lib/api/response"
	"github.com/stepan41k/testMidlware/internal/lib/logger/sl"
	"github.com/stepan41k/testMidlware/internal/storage"
)

type Response struct {
	resp.Response
	Name string `json:"name"`
}

type BookSaver interface {
	SaveBook(domain.Book) (int64, error)
}

func New(log *slog.Logger, bookSaver BookSaver) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.book.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req domain.Book

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		id, err := bookSaver.SaveBook(domain.Book{
			Name: req.Name,
			Author: req.Author,
			Genre: req.Genre,
			Price: req.Price,
		})
		if errors.Is(err, storage.ErrBookExists) {
			log.Info("book already exists", slog.String("name", req.Name))

			render.JSON(w, r, resp.Error("book alread exists"))

			return
		}

		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save book"))

			return
		}

		log.Info("book saved", slog.Int64("id", id))

		render.JSON(w, r, Response{Response: resp.OK(), Name: req.Name})
	}
}