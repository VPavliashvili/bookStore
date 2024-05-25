package books

import (
	"encoding/json"
	"fmt"
	"net/http"
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
//	@Summary      Lists all books
//	@Description  get books
//	@Tags         books
//	@Accept       json
//	@Produce      json
//	@Success      200  {array}   bookDTO
//	@Failure      500  {object}  APIError
//	@Router       /api/store/books [get]
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
//	@Summary      Get book by id
//	@Description  get book
//	@Tags         books
//	@Accept       json
//	@Produce      json
//	@Success      200  {object}   bookDTO
//	@Failure      500  {object}  APIError
//	@Failure      404  {object}  APIError
//	@Router       /api/store/books{id} [get]
func (api API) GetBook(id int, w http.ResponseWriter, r *http.Request) {
	// repoRes, err := api.repo.GetBookById(id)
}
