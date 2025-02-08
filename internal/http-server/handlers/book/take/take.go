package take

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	resp "github.com/stepan41k/testMidlware/internal/lib/api/response"
	"github.com/stepan41k/testMidlware/internal/storage"
)

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	resp.Response
	storage.Book
}

type OneTaker interface {
	TakeOne(name string) (storage.Book, error)
}

func New(log *slog.Logger, oneTaker OneTaker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.book.take.New"

		log = slog.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		
	}
}