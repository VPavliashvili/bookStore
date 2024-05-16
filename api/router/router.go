package router

import (
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

func New() *customMux {

	docs.SwaggerInfo.Title = "Books store API"
	docs.SwaggerInfo.Description = "This is a simple CRUD api implementation for educatinal purposes"
	docs.SwaggerInfo.Version = "1.0"

	router := &customMux{
		ServeMux: http.NewServeMux(),
	}
	router.use(testMiddleware)
	router.handle("/api/info/", infoRouteGroup())

	router.HandleFunc("GET /swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", 6012)),
	))

	return router
}

func infoRouteGroup() http.Handler {
	subRouter := &customMux{
		ServeMux: http.NewServeMux(),
	}

	subRouter.use(testMiddlewareSubRouter)

	subRouter.handleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "health check from myRouter")
	})
	subRouter.handleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "api info from myRouter")
	})
	subRouter.handleFunc("POST /about", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "about as post from myrouter")
	})

	return http.StripPrefix("/api/info", subRouter)
}
