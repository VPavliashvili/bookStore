package books

import (
	"booksapi/api/database"
	"booksapi/logger"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type IBooksRepo interface {
	GetBooks() ([]bookEntity, error)
	GetBookById(int) (bookEntity, error)
	AddBook(bookRequestBody) (int, error)
	RemoveBook(int) error
	UpdateBook(int, bookRequestBody) error
}

type BooksRepo struct{}

func (repo *BooksRepo) GetBooks() ([]bookEntity, error) {
	result := make([]bookEntity, 0)
	query := `SELECT * FROM public.books`

	rows, err := database.Pool.Query(context.Background(), query)
	if err != nil {
		logger.Error(err.Error())
		return result, internalErr{message: err.Error()}
	}

	for rows.Next() {
		var r bookEntity
		err := rows.Scan(&r.ID, &r.Title, &r.Author, &r.Genre,
			&r.NumberOfPages, &r.Price, &r.ReleaseYear)
		if err != nil {
			logger.Error(err.Error())
			return result, internalErr{message: err.Error()}
		}
		result = append(result, r)
	}

	if err = rows.Err(); err != nil {
		logger.Error(err.Error())
		return result, internalErr{message: err.Error()}
	}

	return result, nil
}

func (repo *BooksRepo) GetBookById(id int) (bookEntity, error) {
	query := `SELECT * FROM public.books WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	var b bookEntity
	err := database.Pool.QueryRow(context.Background(), query, args).
		Scan(&b.ID, &b.Title, &b.Author, &b.Genre, &b.NumberOfPages, &b.Price, &b.ReleaseYear)
	if err != nil {
		logger.Error(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return b, notfoundErr{message: err.Error()}
		}
		return b, internalErr{message: err.Error()}
	}

	return b, nil
}

func (repo *BooksRepo) AddBook(b bookRequestBody) (int, error) {
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
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, notfoundErr{message: err.Error()}
		}
		return 0, internalErr{message: err.Error()}
	}

	return id, nil
}

func (repo *BooksRepo) RemoveBook(id int) error {
	query := `DELETE FROM public.books WHERE id = @id`
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := database.Pool.Exec(context.Background(), query, args)
	if err != nil {
		logger.Error(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return notfoundErr{message: err.Error()}
		}
		return internalErr{message: err.Error()}
	}

	return nil
}

func (repo *BooksRepo) UpdateBook(id int, b bookRequestBody) error {
	existing, err := repo.GetBookById(id)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	updated := bookEntity{
		ID:            id,
		Title:         existing.Title,
		Author:        existing.Author,
		Genre:         existing.Genre,
		NumberOfPages: existing.NumberOfPages,
		Price:         existing.Price,
		ReleaseYear:   existing.ReleaseYear,
	}

	if b.Author != nil {
		updated.Author = *b.Author
	}
	if b.Title != nil {
		updated.Title = *b.Author
	}
	if b.Genre != nil {
		updated.Genre = *b.Genre
	}
	if b.NumberOfPages != nil {
		updated.NumberOfPages = b.NumberOfPages
	}
	if b.Price != nil {
		updated.Price = b.Price
	}
	if b.ReleaseYear != nil {
		updated.ReleaseYear = b.ReleaseYear
	}

	logger.Info(fmt.Sprintf("title: %s, author: %s, genre: %s",
		updated.Title, updated.Author, updated.Genre))

	query := `UPDATE public.books
	          SET title = @title, author = @author, genre = @genre, number_of_pages = @number_of_pages,
                            price = @price, release_year = @release_year
              WHERE id = @id`

	args := pgx.NamedArgs{
		"id":              id,
		"title":           updated.Title,
		"author":          updated.Author,
		"genre":           updated.Genre,
		"number_of_pages": updated.NumberOfPages,
		"price":           updated.Price,
		"release_year":    updated.ReleaseYear,
	}

	_, err = database.Pool.Exec(context.Background(), query, args)
	if err != nil {
		logger.Error(err.Error())
		return internalErr{message: err.Error()}
	}

	return nil
}
