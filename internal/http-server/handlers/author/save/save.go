package save

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	resp "github.com/stepan41k/testMidlware/internal/lib/api/response"
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

		
	}
}