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
		w.WriteHeader(http.StatusInternalServerError)
		e := APIError{
			Status:  http.StatusInternalServerError,
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

	w.WriteHeader(http.StatusOK)
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
		w.WriteHeader(http.StatusBadRequest)
		e := APIError{
			Status:  http.StatusBadRequest,
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
			code = http.StatusInternalServerError
		case notfoundErr:
			code = http.StatusNotFound
		case badreqErr:
			code = http.StatusBadRequest
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

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(json[:]))
}

// AddBook adds new books into database
//
//	@Summary		Add new book
//	@Description	adds book
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			newbook	body		bookDTO	true	"request body"
//	@Success		201		{object}	ActionResponse
//	@Failure		500		{object}	APIError
//	@Failure		400		{object}	APIError
//	@Router			/api/books [post]
func (api API) AddBook(w http.ResponseWriter, r *http.Request) {
	var dto bookDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		e := APIError{
			Message: "invalid request model",
			Status:  http.StatusBadRequest,
		}
		fmt.Fprint(w, e.Error())
		return
	}

	if len(dto.Title) == 0 || len(dto.Author) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		e := APIError{
			Message: "required fields are not set, won't save the data",
			Status:  http.StatusBadRequest,
		}
		fmt.Fprint(w, e.Error())
		return
	}

	entitiy := dto.ToEntity()
	id, err := api.repo.AddBook(entitiy)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		e := APIError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
		fmt.Fprint(w, e.Error())
		return
	}

	resp := ActionResponse{
		ResourceId: id,
	}
	j, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(j[:]))
}
