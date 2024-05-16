package middlewares

import (
	"fmt"
	"net/http"
)

func testMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("http request method %s, url %s\n", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func testMiddlewareSubRouter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("on router %s middleware run from subrouter\n", r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

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
