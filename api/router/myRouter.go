package router

import (
	"fmt"
	"net/http"
)

type myRouter struct {
	*http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

func (m *myRouter) handle(pattern string, handler http.Handler) {
	for _, middleware := range m.middlewares {
		handler = middleware(handler)
	}
	m.Handle(pattern, handler)
}

func (m *myRouter) handleFunc(pattern string, handler http.HandlerFunc) {
	m.handle(pattern, handler)
}

func (m *myRouter) use(middlewares ...func(http.Handler) http.Handler) *myRouter {
	m.middlewares = append(m.middlewares, middlewares...)
	return m
}

func NewMyRouter() *myRouter {
	router := &myRouter{
		ServeMux: http.NewServeMux(),
	}
	router.use(testMiddleware)
	router.handle("/api/info/", info())

	// mux.HandleFunc("GET /swagger/*", httpSwagger.Handler(
	// 	httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", app.config.port)),
	// ))

	return router
}

func info() http.Handler {
	subRouter := &myRouter{
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
