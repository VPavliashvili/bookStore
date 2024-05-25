package books

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type API struct {
	repo IBooksRepo
}

func New() API {
	return API{
		repo: &BooksRepo{},
	}
}

// GetBooks returns all books
//
//	@Summary		Lists all books
//	@Description	get books
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		bookDTO
//	@Failure		500	{object}	APIError
//	@Router			/api/books [get]
func (api API) GetBooks(w http.ResponseWriter, r *http.Request) {
	repoRes, err := api.repo.GetBooks()
	if err != nil {
		w.WriteHeader(500)
		e := APIError{
			Status:  500,
			Message: err.Error(),
		}
		fmt.Fprint(w, e.Error())
		return
	}

	dtos := make([]bookDTO, 0)
	for _, b := range repoRes {
		dtos = append(dtos, b.ToDto())
	}

	json, _ := json.Marshal(dtos)

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", string(json[:]))
}

// GetBook returns a book by id
//
//	@Summary		Get book by id
//	@Description	get book
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Book ID"
//	@Success		200	{object}	bookDTO
//	@Failure		500	{object}	APIError
//	@Failure		400	{object}	APIError
//	@Failure		404	{object}	APIError
//	@Router			/api/books/{id} [get]
func (api API) GetBook(w http.ResponseWriter, r *http.Request) {
	p := r.PathValue("id")
	id, err := strconv.Atoi(p)
	if err != nil {
		w.WriteHeader(400)
		e := APIError{
			Status:  400,
			Message: "only accept integer values as {id} path parameter",
		}
		fmt.Fprint(w, e.Error())
		return
	}

	book, err := api.repo.GetBookById(id)
	if err != nil {
		var code int
		switch err.(type) {
		case internalErr:
			code = 500
		case notfoundErr:
			code = 404
		case badreqErr:
			code = 400
		}
		w.WriteHeader(code)
		e := APIError{
			Status:  code,
			Message: err.Error(),
		}
		fmt.Fprint(w, e.Error())
        return
	}

	dto := book.ToDto()
	json, _ := json.Marshal(dto)

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", string(json[:]))
}
