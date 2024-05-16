package books

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type API struct {
	repo IBooksRepo
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
		json, _ := json.Marshal(e)
		fmt.Fprint(w, string(json[:]))
		return
	}

	var dtos []bookDTO
	for _, b := range repoRes {
		dtos = append(dtos, b.ToDto())
	}

	json, _ := json.Marshal(dtos)

	fmt.Fprintf(w, "%s\n", string(json[:]))
}

func New() API {
	return API{
		repo: &BooksRepo{},
	}
}
