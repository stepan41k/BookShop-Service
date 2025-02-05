package save

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	resp "github.com/stepan41k/testMidlware/internal/lib/api/response"
)

type Request struct {
	Name string `json:"name"`
	Author string `json:"author"`
	Genre string `json:"genre"`
	Price float32 `json:"float"`
}

type Response struct {
	resp.Response
	Name string `json:"name"`
}

type BookSaver interface {
	SaveBook(name string, author string, genre string, price float32) (int64, error)
}

func New(log *slog.Logger, bookSaver BookSaver) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.book.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
	}
}