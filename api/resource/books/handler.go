package books

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func writeAPIErr(err APIError, w http.ResponseWriter) {
	w.WriteHeader(err.Status)
	fmt.Fprint(w, err.Error())
}

func writeErr(err error, status int, w http.ResponseWriter) {
	e := APIError{
		Status:  status,
		Message: err.Error(),
	}

	writeAPIErr(e, w)
}

func getRepoErrcode(err error) int {
	var code int
	switch err.(type) {
	case internalErr:
		code = http.StatusInternalServerError
	case notfoundErr:
		code = http.StatusNotFound
	}
	return code
}

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
		writeErr(err, http.StatusInternalServerError, w)
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
		e := APIError{
			Status:  http.StatusBadRequest,
			Message: "only accept integer values as {id} path parameter",
		}
		writeAPIErr(e, w)
		return
	}

	book, err := api.repo.GetBookById(id)
	if err != nil {
		code := getRepoErrcode(err)
		writeErr(err, code, w)
		return
	}

	dto := book.ToDto()
	json, _ := json.Marshal(dto)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(json[:]))
}

// AddBook adds new book into database
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
	defer r.Body.Close()
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&dto)
	if err != nil {
		e := APIError{
			Message: "invalid request model",
			Status:  http.StatusBadRequest,
		}
		writeAPIErr(e, w)
		return
	}

	if len(dto.Title) == 0 || len(dto.Author) == 0 {
		e := APIError{
			Message: "required fields are not set, won't save the data",
			Status:  http.StatusBadRequest,
		}
		writeAPIErr(e, w)
		return
	}

	entitiy := dto.ToEntity()
	id, err := api.repo.AddBook(entitiy)

	if err != nil {
		writeErr(err, http.StatusInternalServerError, w)
		return
	}

	resp := ActionResponse{
		ResourceId: id,
	}
	j, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, string(j[:]))
}

// DeleteBook deletes existing book from database
//
//	@Summary		Remove book record
//	@Description	removes book record
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"book record Id"
//	@Success		204
//	@Failure		500	{object}	APIError
//	@Failure		400	{object}	APIError
//	@Failure		404	{object}	APIError
//	@Router			/api/books/{id} [delete]
func (api API) RemoveBook(w http.ResponseWriter, r *http.Request) {
	p := r.PathValue("id")
	id, err := strconv.Atoi(p)
	if err != nil {
		e := APIError{
			Status:  http.StatusBadRequest,
			Message: "only accept integer values as {id} path parameter",
		}
		writeAPIErr(e, w)
		return
	}

	err = api.repo.RemoveBook(id)
	if err != nil {
		code := getRepoErrcode(err)
		writeErr(err, code, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w, "")

}

// UpdateBook updates existing book from database
//
//	@Summary		Update book record
//	@Description	Updates book record
//	@Tags			books
//	@Accept			json
//	@Produce		json
//	@Param			id		path	int		true	"book record Id"
//	@Param			book	body	bookDTO	true	"request body"
//	@Success		200
//	@Failure		500	{object}	APIError
//	@Failure		400	{object}	APIError
//	@Failure		404	{object}	APIError
//	@Router			/api/books/{id} [patch]
func (api API) UpdateBook(w http.ResponseWriter, r *http.Request) {
	p := r.PathValue("id")
	id, err := strconv.Atoi(p)
	if err != nil {
		e := APIError{
			Status:  http.StatusBadRequest,
			Message: "only accept integer values as {id} path parameter",
		}
		writeAPIErr(e, w)
		return
	}

	var dto bookDTO
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&dto)
	if err != nil {
		e := APIError{
			Message: "invalid request model",
			Status:  http.StatusBadRequest,
		}
		writeAPIErr(e, w)
		return
	}

	entitiy := dto.ToEntity()
	err = api.repo.UpdateBook(id, entitiy)
	if err != nil {
		code := getRepoErrcode(err)
		writeErr(err, code, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Print(w, "")

}
