package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/stepan41k/testMidlware/internal/lib/api/response"
	"github.com/stepan41k/testMidlware/internal/lib/logger/sl"
	"github.com/stepan41k/testMidlware/internal/storage"
)

type Request struct {
	Genre string `json:"genre" validate:"required"`
}

type Response struct {
	resp.Response
	Genre string `json:"genre"`
}

type GenreSaver interface {
	SaveGenre(genre string) (int64, error)
}

func New(log *slog.Logger, genreSaver GenreSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.genre.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}

		if err != nil {
			log.Error("failed to decode request body")

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

		genre := req.Genre
		
		id, err := genreSaver.SaveGenre(req.Genre)
		if errors.Is(err, storage.ErrGenreExists) {
			log.Info("genre already exists", slog.String("genre", req.Genre))

			render.JSON(w, r, resp.Error("genre already exists"))

			return
		}

		if err != nil {
			log.Error("failed to add genre", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add genre"))

			return
		}

		log.Info("genre added", slog.Int64("id", id))

		responseOK(w, r, genre)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, genre string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Genre: genre,
	})
}