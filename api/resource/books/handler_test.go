package books

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type fakeWriter struct {
	input        string
	headerStatus int
}

func (w fakeWriter) Header() http.Header {
	panic("unimplemented")
}

func (w *fakeWriter) Write(p []byte) (int, error) {
	w.input = string(p[:])
	return 0, nil
}

func (w *fakeWriter) WriteHeader(statusCode int) {
	w.headerStatus = statusCode
}

type fakeRepo struct {
	pluralReturner   func() ([]bookEntity, error)
	singleReturner   func(int) (bookEntity, error)
	addbookAction    func(bookEntity) (int, error)
	removeBookAction func(int) error
	updateBookAction func(int, bookEntity) error
}

func (r fakeRepo) GetBookById(id int) (bookEntity, error) {
	return r.singleReturner(id)
}

func (r fakeRepo) GetBooks() ([]bookEntity, error) {
	return r.pluralReturner()
}

func (r fakeRepo) AddBook(e bookEntity) (int, error) {
	return r.addbookAction(e)
}

func (r fakeRepo) RemoveBook(id int) error {
	return r.removeBookAction(id)
}

func (r fakeRepo) UpdateBook(id int, b bookEntity) error {
	return r.updateBookAction(id, b)
}

func TestGetBooks(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		w        *fakeWriter
		expected struct {
			data         string
			headerStatus int
		}
	}{
		{
			repo: fakeRepo{pluralReturner: func() ([]bookEntity, error) {
				return nil, errors.New("fake err")
			}},
			w: &fakeWriter{},
			expected: struct {
				data         string
				headerStatus int
			}{
				data: APIError{
					Status:  http.StatusInternalServerError,
					Message: "fake err",
				}.Error(),
				headerStatus: http.StatusInternalServerError,
			},
		},
		{
			repo: fakeRepo{
				pluralReturner: func() ([]bookEntity, error) {
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
				},
			},
			w: &fakeWriter{},
			expected: struct {
				data         string
				headerStatus int
			}{
				data: func() string {
					dtos := []bookDTO{
						{
							Title:         "The Fellowship of the Ring",
							Author:        "JRR Tolkien",
							Price:         20,
							NumberOfPages: 432,
							Genre:         "fantasy",
							ReleaseYear:   1954,
						},
					}
					json, _ := json.Marshal(dtos)
					return string(json[:])
				}(),
				headerStatus: http.StatusOK,
			},
		},
		{
			repo: fakeRepo{
				pluralReturner: func() ([]bookEntity, error) {
					return []bookEntity{}, nil
				},
			},
			w: &fakeWriter{},
			expected: struct {
				data         string
				headerStatus int
			}{
				data:         "[]",
				headerStatus: http.StatusOK,
			},
		},
		{
			repo: fakeRepo{
				pluralReturner: func() ([]bookEntity, error) {
					return nil, nil
				},
			},
			w: &fakeWriter{},
			expected: struct {
				data         string
				headerStatus int
			}{
				data:         "[]",
				headerStatus: http.StatusOK,
			},
		},
	}

	for _, tc := range tcases {
		api := API{repo: tc.repo}
		api.GetBooks(tc.w, nil)
		if tc.expected.data != tc.w.input {
			t.Errorf("GetBooks failed\nexpected %v\ngot %s", tc.expected.data, tc.w.input)
		}
		if tc.expected.headerStatus != tc.w.headerStatus {
			t.Errorf("GetBooks response header failed\nexpected %v\ngot  %v",
				tc.expected.headerStatus, tc.w.headerStatus)
		}
	}
}

func TestGetBookById(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		w        *fakeWriter
		req      *http.Request
		expected struct {
			data         string
			headerStatus int
		}
	}{
		{
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("id", "wrongStr")

				return rq
			}(),
			repo: fakeRepo{},
			w:    &fakeWriter{},
			expected: struct {
				data         string
				headerStatus int
			}{
				data: APIError{
					Status:  http.StatusBadRequest,
					Message: "only accept integer values as {id} path parameter",
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			repo: fakeRepo{singleReturner: func(i int) (bookEntity, error) {
				return bookEntity{}, notfoundErr{fmt.Sprintf("resource not found err, at id -> %v", i)}
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("id", "123")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: APIError{
					Status:  http.StatusNotFound,
					Message: "resource not found err, at id -> 123",
				}.Error(),
				headerStatus: http.StatusNotFound,
			},
		},
		{
			repo: fakeRepo{singleReturner: func(i int) (bookEntity, error) {
				return bookEntity{}, internalErr{message: "internal err"}
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("id", "10")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: APIError{
					Status:  http.StatusInternalServerError,
					Message: "internal err",
				}.Error(),
				headerStatus: http.StatusInternalServerError,
			},
		},
		{
			repo: fakeRepo{singleReturner: func(i int) (bookEntity, error) {
				return bookEntity{Title: "Superman Red Son"}, nil
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("id", "10")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: func() string {
					dto := bookDTO{
						Title: "Superman Red Son",
					}
					json, _ := json.Marshal(dto)
					return string(json[:])
				}(),
				headerStatus: http.StatusOK,
			},
		},
	}

	for _, tc := range tcases {
		api := API{repo: tc.repo}
		api.GetBook(tc.w, tc.req)
		if tc.expected.data != tc.w.input {
			t.Errorf("GetBook failed\nexpected %v\ngot %s", tc.expected.data, tc.w.input)
		}
		if tc.expected.headerStatus != tc.w.headerStatus {
			t.Errorf("GetBook response header failed\nexpected %v\ngot  %v",
				tc.expected.headerStatus, tc.w.headerStatus)
		}
	}
}

func TestAddBook(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		w        *fakeWriter
		req      *http.Request
		expected struct {
			data         string
			headerStatus int
		}
	}{
		{
			repo: fakeRepo{
				addbookAction: func(bookEntity) (int, error) {
					return 0, badreqErr{message: "required fields are not set, won't save the data"}
				},
			},
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				j, _ := json.Marshal(bookDTO{})
				rq.Body = io.NopCloser(strings.NewReader(string(j[:])))

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: func() string {
					e := APIError{
						Message: "required fields are not set, won't save the data",
						Status:  http.StatusBadRequest,
					}
					return e.Error()
				}(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			repo: fakeRepo{},
			w:    &fakeWriter{},
			req: func() *http.Request {
				msg := `{"tst":"value"}`
				rq, _ := http.NewRequest("POST", "", strings.NewReader(msg))

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: func() string {
					e := APIError{
						Message: "invalid request model",
						Status:  http.StatusBadRequest,
					}
					return e.Error()
				}(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			repo: fakeRepo{addbookAction: func(e bookEntity) (int, error) {
				if e.Title != "test" || e.Author != "tst" || e.Genre != "idk" ||
					e.NumberOfPages != 1 || e.Price != 2 || e.ReleaseYear != 3 {
					return 0, errors.New("")
				}
				return 1, nil
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				d := bookDTO{
					Title:         "test",
					Author:        "tst",
					Genre:         "idk",
					NumberOfPages: 1,
					Price:         2,
					ReleaseYear:   3,
				}
				j, _ := json.Marshal(d)
				rq, _ := http.NewRequest("POST", "", strings.NewReader(string(j[:])))

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: func() string {
					a := ActionResponse{ResourceId: 1}
					j, _ := json.Marshal(a)
					return string(j[:])
				}(),
				headerStatus: http.StatusCreated,
			},
		},
		{
			repo: fakeRepo{addbookAction: func(e bookEntity) (int, error) {
				return 0, internalErr{message: "internal error"}
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				d := bookDTO{
					Title:  "test",
					Author: "tst",
				}
				j, _ := json.Marshal(d)
				rq, _ := http.NewRequest("POST", "", strings.NewReader(string(j[:])))

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: func() string {
					a := APIError{Message: "internal error", Status: http.StatusInternalServerError}
					j, _ := json.Marshal(a)
					return string(j[:])
				}(),
				headerStatus: http.StatusInternalServerError,
			},
		}}

	for _, tc := range tcases {
		api := API{repo: tc.repo}
		api.AddBook(tc.w, tc.req)
		if tc.expected.data != tc.w.input {
			t.Errorf("AddBook failed\nexpected %v\ngot %s", tc.expected.data, tc.w.input)
		}
		if tc.expected.headerStatus != tc.w.headerStatus {
			t.Errorf("AddBook response header failed\nexpected %v\ngot  %v",
				tc.expected.headerStatus, tc.w.headerStatus)
		}
	}
}

func TestRemoveBookById(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		w        *fakeWriter
		req      *http.Request
		expected struct {
			data         string
			headerStatus int
		}
	}{
		{
			repo: fakeRepo{},
			w:    &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("id", "wrongStr")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: APIError{
					Status:  http.StatusBadRequest,
					Message: "only accept integer values as {id} path parameter",
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			repo: fakeRepo{removeBookAction: func(i int) error {
				return notfoundErr{fmt.Sprintf("resource not found err, at id -> %v", i)}
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("id", "123")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: APIError{
					Status:  http.StatusNotFound,
					Message: "resource not found err, at id -> 123",
				}.Error(),
				headerStatus: http.StatusNotFound,
			},
		},
		{
			repo: fakeRepo{removeBookAction: func(i int) error {
				return internalErr{message: "internal err"}
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("id", "10")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: APIError{
					Status:  http.StatusInternalServerError,
					Message: "internal err",
				}.Error(),
				headerStatus: http.StatusInternalServerError,
			},
		},
		{
			repo: fakeRepo{removeBookAction: func(i int) error {
				return nil
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("id", "10")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data:         "",
				headerStatus: http.StatusNoContent,
			},
		},
	}

	for _, tc := range tcases {
		api := API{repo: tc.repo}
		api.RemoveBook(tc.w, tc.req)
		if tc.expected.data != tc.w.input {
			t.Errorf("GetBook failed\nexpected %v\ngot %s", tc.expected.data, tc.w.input)
		}
		if tc.expected.headerStatus != tc.w.headerStatus {
			t.Errorf("GetBook response header failed\nexpected %v\ngot  %v",
				tc.expected.headerStatus, tc.w.headerStatus)
		}
	}
}

func TestUpdateBookById(t *testing.T) {
	tcases := []struct {
		repo     fakeRepo
		w        *fakeWriter
		req      *http.Request
		expected struct {
			data         string
			headerStatus int
		}
	}{
		{
			req: func() *http.Request {
				rq := &http.Request{}
				rq.SetPathValue("id", "wrongStr")

				return rq
			}(),
			repo: fakeRepo{},
			w:    &fakeWriter{},
			expected: struct {
				data         string
				headerStatus int
			}{
				data: APIError{
					Status:  http.StatusBadRequest,
					Message: "only accept integer values as {id} path parameter",
				}.Error(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			repo: fakeRepo{},
			w:    &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				j := `{"tst":"tst"}`
				rq.Body = io.NopCloser(strings.NewReader(string(j[:])))
				rq.SetPathValue("id", "10")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: func() string {
					e := APIError{
						Message: "invalid request model",
						Status:  http.StatusBadRequest,
					}
					return e.Error()
				}(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			repo: fakeRepo{},
			w:    &fakeWriter{},
			req: func() *http.Request {
				rq := &http.Request{}
				j := `random text`
				rq.Body = io.NopCloser(strings.NewReader(string(j[:])))
				rq.SetPathValue("id", "10")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: func() string {
					e := APIError{
						Message: "invalid request model",
						Status:  http.StatusBadRequest,
					}
					return e.Error()
				}(),
				headerStatus: http.StatusBadRequest,
			},
		},
		{
			repo: fakeRepo{updateBookAction: func(i int, b bookEntity) error {
				if b.Title != "test" || b.Author != "tst" || b.Genre != "idk" ||
					b.NumberOfPages != 1 || b.Price != 2 || b.ReleaseYear != 3 {
					return errors.New("")
				}
				return nil
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				d := bookDTO{
					Title:         "test",
					Author:        "tst",
					Genre:         "idk",
					NumberOfPages: 1,
					Price:         2,
					ReleaseYear:   3}
				j, _ := json.Marshal(d)
				rq, _ := http.NewRequest("POST", "", strings.NewReader(string(j[:])))
				rq.SetPathValue("id", "10")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data:         "",
				headerStatus: http.StatusOK,
			},
		},
		{
			repo: fakeRepo{updateBookAction: func(i int, b bookEntity) error {
				return internalErr{"internal error"}
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				d := bookDTO{}
				j, _ := json.Marshal(d)
				rq, _ := http.NewRequest("POST", "", strings.NewReader(string(j[:])))
				rq.SetPathValue("id", "10")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: APIError{
					Status:  http.StatusInternalServerError,
					Message: "internal error",
				}.Error(),
				headerStatus: http.StatusInternalServerError,
			},
		},
		{
			repo: fakeRepo{updateBookAction: func(i int, b bookEntity) error {
				return notfoundErr{"not found error"}
			}},
			w: &fakeWriter{},
			req: func() *http.Request {
				d := bookDTO{}
				j, _ := json.Marshal(d)
				rq, _ := http.NewRequest("POST", "", strings.NewReader(string(j[:])))
				rq.SetPathValue("id", "10")

				return rq
			}(),
			expected: struct {
				data         string
				headerStatus int
			}{
				data: APIError{
					Status:  http.StatusNotFound,
					Message: "not found error",
				}.Error(),
				headerStatus: http.StatusNotFound,
			},
		},
	}

	for _, tc := range tcases {
		api := API{repo: tc.repo}
		api.UpdateBook(tc.w, tc.req)
		if tc.expected.data != tc.w.input {
			t.Errorf("GetBook failed\nexpected %v\ngot %s", tc.expected.data, tc.w.input)
		}
		if tc.expected.headerStatus != tc.w.headerStatus {
			t.Errorf("GetBook response header failed\nexpected %v\ngot  %v",
				tc.expected.headerStatus, tc.w.headerStatus)
		}
	}
}
