package middlewares

import (
	"booksapi/api/router"
	"booksapi/logger"
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

		next.ServeHTTP(w, r)
	})
}

func LogRequestResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rww := router.NewResponseWriterWrapper(w)
		defer func() {
			msg := logger.GetRequestResponseLog(rww, r)
			logger.Info(msg)
		}()

		next.ServeHTTP(rww, r)
	})
}
