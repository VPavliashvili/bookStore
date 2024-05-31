package router

import (
	"bytes"
	"net/http"
)

// ResponseWriterWrapper struct is used for logging middleware
type ResponseWriterWrapper struct {
	W          *http.ResponseWriter
	Body       *bytes.Buffer
	StatusCode *int
}

func NewResponseWriterWrapper(w http.ResponseWriter) ResponseWriterWrapper {
	var buf bytes.Buffer
	var statusCode = 200
	return ResponseWriterWrapper{
		W:          &w,
		Body:       &buf,
		StatusCode: &statusCode,
	}
}

// overwrites Write() function
func (rww ResponseWriterWrapper) Write(buf []byte) (int, error) {
	rww.Body.Write(buf)
	return (*rww.W).Write(buf)
}

// overwrites Header() function
func (rww ResponseWriterWrapper) Header() http.Header {
	return (*rww.W).Header()
}

// overwrites WriteHeader() function
func (rww ResponseWriterWrapper) WriteHeader(statusCode int) {
	(*rww.StatusCode) = statusCode
	(*rww.W).WriteHeader(statusCode)
}
