package router

import (
	system "booksapi/api/resource/system"
	"booksapi/api/router/middlewares"
	"booksapi/docs"
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type customMux struct {
	*http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

func (m *customMux) handle(pattern string, handler http.Handler) {
	for _, middleware := range m.middlewares {
		handler = middleware(handler)
	}
	m.Handle(pattern, handler)
}

func (m *customMux) handleFunc(pattern string, handler http.HandlerFunc) {
	m.handle(pattern, handler)
}

func (m *customMux) use(middlewares ...func(http.Handler) http.Handler) *customMux {
	m.middlewares = append(m.middlewares, middlewares...)
	return m
}

func New(systemApi system.API) *customMux {

	docs.SwaggerInfo.Title = "Books store API"
	docs.SwaggerInfo.Description = "This is a simple CRUD api implementation for educatinal purposes"
	docs.SwaggerInfo.Version = "1.0"

	router := &customMux{
		ServeMux: http.NewServeMux(),
	}
	router.use(middlewares.ContentTypeJSON)
	router.handle("/api/system/", systemRouteGroup(systemApi))

	router.HandleFunc("GET /swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", 6012)),
	))

	return router
}

func systemRouteGroup(api system.API) http.Handler {
	subRouter := &customMux{
		ServeMux: http.NewServeMux(),
	}

	subRouter.handleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		api.HandleHealth(w, r)
	})
	subRouter.handleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		api.HandleAbout(w, r)
	})

	return http.StripPrefix("/api/system", subRouter)
}
