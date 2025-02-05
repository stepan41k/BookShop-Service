package repository

import (
	"testing"
	"github.com/stepan41k/testMidlware/internal/pkg/models"
)

const connStr = "postgres://postgres:admin@localhost:5432/cources"

func TestBooksCRUD(t *testing.T) {

	db, err := New(connStr)
	if err != nil {
		t.Fatal(err.Error())
	}

	_, err = db.GetBooks()
	if err != nil {
		t.Fatal(err.Error())
	}

	book := models.Book{GenreID: 1, AuthorID: 1, Name: "Идиот", Price: 150}
	id, err := db.NewBook(book)
	if err != nil || id == 0{
		t.Fatal(err.Error())
	}

	_, err = db.GetBookByID(id)
	if err != nil || id == 0{
		t.Fatal(err.Error())
	}
}