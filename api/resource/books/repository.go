package books

import (
	"booksapi/api/database"
	"booksapi/logger"
	"context"

	"github.com/jackc/pgx/v5"
)

type IBooksRepo interface {
	GetBooks() ([]bookEntity, error)
	GetBookById(int) (bookEntity, error)
	AddBook(bookEntity) (int, error)
	RemoveBook(int) error
	UpdateBook(int, bookEntity) error
}

type BooksRepo struct{}

func (repo *BooksRepo) GetBooks() ([]bookEntity, error) {
	return []bookEntity{
		{
			Title:         "The Fellowship of the Ring",
			Author:        "JRR Tolkien",
			Price:         20,
			NumberOfPages: 432,
			Genre:         "fantasy",
			ReleaseYear:   1954,
		},
	}, nil
}

func (repo *BooksRepo) GetBookById(id int) (bookEntity, error) {
	return bookEntity{
		Title:         "The Two Towers",
		Author:        "JRR Tolkien",
		Genre:         "fantasy",
		NumberOfPages: 352,
		Price:         20,
		ReleaseYear:   1954,
	}, nil
}

func (repo *BooksRepo) AddBook(b bookEntity) (int, error) {
	query := `INSERT INTO public.books
                (title, author, genre, number_of_pages, price, release_year)
                VALUES(@title, @author, @genre, @number_of_pages, @price, @release_year) RETURNING id`
	args := pgx.NamedArgs{
		"title":           b.Title,
		"author":          b.Author,
		"genre":           b.Genre,
		"number_of_pages": b.NumberOfPages,
		"price":           b.Price,
		"release_year":    b.ReleaseYear,
	}

	var id int
	err := database.Pool.QueryRow(context.Background(), query, args).Scan(&id)
	if err != nil {
		logger.Error(err.Error())
		return 0, internalErr{message: err.Error()}
	}

	return id, nil
}

func (repo *BooksRepo) RemoveBook(int) error {
	return nil
}

func (repo *BooksRepo) UpdateBook(int, bookEntity) error {
	return nil
}
