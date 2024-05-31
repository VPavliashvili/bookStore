package middlewares

import (
	"booksapi/logger"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

func ContentTypeJSON(next http.Handler) http.Handler {
	const (
		HeaderKeyContentType       = "Content-Type"
		HeaderValueContentTypeJSON = "application/json;charset=utf8"
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderKeyContentType, HeaderValueContentTypeJSON)
		next.ServeHTTP(w, r)
	})
}

func RequestID(next http.Handler) http.Handler {
	const (
		XRequestIDKey = "XRequestID"
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xRequestID := uuid.NewString()
        w.Header().Set(XRequestIDKey, xRequestID)

		logger.SetNewUUID(xRequestID)

		next.ServeHTTP(w, r)
	})
}

func LogRequestResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := getRequestLog(r)

		rww := newResponseWriterWrapper(w)

		defer func() {
			resp := getResponseLog(rww)

			msg := fmt.Sprintf("[Request: %s] [Response: %s]\n", req, resp)
			logger.Info(msg)
		}()

		next.ServeHTTP(rww, r)
	})
}

// responseWriterWrapper struct is used to log the response
type responseWriterWrapper struct {
	w          *http.ResponseWriter
	body       *bytes.Buffer
	statusCode *int
}

func newResponseWriterWrapper(w http.ResponseWriter) responseWriterWrapper {
	var buf bytes.Buffer
	var statusCode = 200
	return responseWriterWrapper{
		w:          &w,
		body:       &buf,
		statusCode: &statusCode,
	}
}

// overwrites Write() function
func (rww responseWriterWrapper) Write(buf []byte) (int, error) {
	rww.body.Write(buf)
	return (*rww.w).Write(buf)
}

// overwrites Header() function
func (rww responseWriterWrapper) Header() http.Header {
	return (*rww.w).Header()
}

// overwrites WriteHeader() function
func (rww responseWriterWrapper) WriteHeader(statusCode int) {
	(*rww.statusCode) = statusCode
	(*rww.w).WriteHeader(statusCode)
}

func getResponseLog(rww responseWriterWrapper) string {
	var buf bytes.Buffer

	buf.WriteString("Response:")

	buf.WriteString("Headers:")
	for k, v := range (*rww.w).Header() {
		buf.WriteString(fmt.Sprintf("%s: %v", k, v))
	}

	buf.WriteString(fmt.Sprintf(" Status Code: %d ", *(rww.statusCode)))

	buf.WriteString("Body ")
	buf.WriteString(rww.body.String())
	return buf.String()
}

func getRequestLog(r *http.Request) string {
	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	data := string(body[:])

	url := r.URL.String()
	endpoint := fmt.Sprintf("%s %s", r.Method, url)

	params := r.URL.Query().Encode()

	log := fmt.Sprintf("endpoint: %s, params: %s, body: %s", endpoint, params, data)
	return log
}
