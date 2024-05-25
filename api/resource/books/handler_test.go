package books

import (
	"encoding/json"
	"errors"
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
	returner func() ([]bookEntity, error)
}

func (r fakeRepo) GetBookById(int) (bookEntity, error) {
	panic("unimplemented")
}

func (r fakeRepo) GetBooks() ([]bookEntity, error) {
	return r.returner()
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
			repo: fakeRepo{returner: func() ([]bookEntity, error) {
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
				returner: func() ([]bookEntity, error) {
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
				returner: func() ([]bookEntity, error) {
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
				returner: func() ([]bookEntity, error) {
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
