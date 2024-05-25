package books

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
	pluralReturner func() ([]bookEntity, error)
	singleReturner func(int) (bookEntity, error)
}

func (r fakeRepo) GetBookById(id int) (bookEntity, error) {
	return r.singleReturner(id)
}

func (r fakeRepo) GetBooks() ([]bookEntity, error) {
	return r.pluralReturner()
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
					Status:  500,
					Message: "fake err",
				}.Error(),
				headerStatus: 500,
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
				headerStatus: 200,
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
				headerStatus: 200,
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
				headerStatus: 200,
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
					Status:  400,
					Message: "only accept integer values as {id} path parameter",
				}.Error(),
				headerStatus: 400,
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
					Status:  404,
					Message: "resource not found err, at id -> 123",
				}.Error(),
				headerStatus: 404,
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
					Status:  500,
					Message: "internal err",
				}.Error(),
				headerStatus: 500,
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
				headerStatus: 200,
			},
		},
	}

	for _, tc := range tcases {
		api := API{repo: tc.repo}
		api.GetBook(tc.w, tc.req)
		if tc.expected.data != tc.w.input {
			t.Errorf("GetBooks failed\nexpected %v\ngot %s", tc.expected.data, tc.w.input)
		}
		if tc.expected.headerStatus != tc.w.headerStatus {
			t.Errorf("GetBooks response header failed\nexpected %v\ngot  %v",
				tc.expected.headerStatus, tc.w.headerStatus)
		}
	}
}
