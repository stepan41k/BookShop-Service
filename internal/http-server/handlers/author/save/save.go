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
	Author string `json:"author"`
}

type Response struct {
	resp.Response
	Author string `json:"author"`
}

type AuthorSaver interface {
	SaveAuthor(author string) (int64, error)
}

func New(log *slog.Logger, authorSaver AuthorSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.author.save.New"

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
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request body"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		author := req.Author
		
		id, err := authorSaver.SaveAuthor(author)
		if errors.Is(err, storage.ErrAuthorExists) {
			log.Info("author already exists", slog.String("author", req.Author))

			render.JSON(w, r, resp.Error("author already exists"))

			return
		}

		if err != nil {
			log.Error("failed to save author")

			render.JSON(w, r, resp.Error("failed to save author"))

			return
		}

		log.Info("author saved", slog.Int64("id", id))

		responseOK(w, r, author)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, author string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		Author: author,
	})
}