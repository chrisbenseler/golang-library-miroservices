package domain

import (
	"database/sql"
	"librarymanager/books/common"
	"math/rand"
	"time"
)

//Repository book repository (persistence)
type Repository interface {
	Save(title string, year int, createdByID string) (*Book, error)
	Get(id string) (*Book, common.Error)
	All() (*[]Book, error)
	Destroy(id string) error
}

type repositoryStruct struct {
	db *sql.DB
}

//NewBookRepository create a new book repository
func NewBookRepository(database *sql.DB) Repository {

	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS book (id STRING PRIMARY KEY, title TEXT, year INTEGER, createdByID TEXT)")
	statement.Exec()

	return &repositoryStruct{
		db: database,
	}
}

//Save book
func (r *repositoryStruct) Save(title string, year int, createdByID string) (*Book, error) {

	book := NewBook(GenerateID(), title, year, createdByID)

	statement, _ := r.db.Prepare("INSERT INTO book (id, title, year, createdByID) VALUES (?, ?, ?, ?)")

	if _, err := statement.Exec(book.ID, book.Title, book.Year, book.CreatedByID); err != nil {
		return nil, err
	}

	return book, nil
}

//Get get a book by its id
func (r *repositoryStruct) Get(id string) (*Book, common.Error) {

	book := &Book{}

	rows, err := r.db.Query("SELECT 1 title, year, createdByID FROM book WHERE id = '" + id + "' LIMIT 1")

	if err != nil {
		return nil, common.NewBadRequestError("No book found for the given ID")
	}

	for rows.Next() {

		var title string
		var year int
		var createdByID string
		rows.Scan(&title, &year, &createdByID)
		book = NewBook(id, title, year, createdByID)

	}

	if book.ID == "" {
		return nil, common.NewNotFoundError("No book found for the given ID")
	}

	return book, nil

}

//All list all books
func (r *repositoryStruct) All() (*[]Book, error) {

	books := []Book{}

	rows, _ := r.db.Query("SELECT id, title, year, createdByID FROM book")

	for rows.Next() {
		var id string
		var title string
		var year int
		var createdByID string
		rows.Scan(&id, &title, &year, &createdByID)
		book := NewBook(id, title, year, createdByID)
		books = append(books, *book)
	}

	return &books, nil

}

//Destroy destroy a book by its id
func (r *repositoryStruct) Destroy(id string) error {

	statement, _ := r.db.Prepare("DELETE FROM book WHERE id = ?")

	_, err := statement.Exec(id)
	return err

}

//GenerateID method
func GenerateID() string {

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)

}
