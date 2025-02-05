package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stepan41k/testMidlware/internal/pkg/models"
)

func (api *api) books(w http.ResponseWriter, r *http.Request) {

	switch r.Method{
	case http.MethodGet:

		vars := mux.Vars(r)
		stringID, ok := vars["id"]
		if ok{
			id, err := strconv.Atoi(stringID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			data, err := api.db.GetBookByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			err = json.NewEncoder(w).Encode(data)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			return

		}

		data, err := api.db.GetBooks()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPost:
		var book models.Book
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, err := api.db.NewBook(book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodDelete:
		vars := mux.Vars(r)
		name, ok := vars["name"]
		if !ok{
			http.Error(w, "error in parsing params", http.StatusInternalServerError)
			return
		}

		err := api.db.DeleteBook(name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	}
	
}

