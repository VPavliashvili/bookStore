package router

import (
	"net/http"
)

type Group struct {
	*CustomMux
	pattern string
	action  func(this *Group) http.Handler
}

type CustomMux struct {
	*http.ServeMux
	middlewares []func(http.Handler) http.Handler
	groups      []*Group
}

func (m *CustomMux) HandleRoute(pattern string, handler http.Handler) {
	for _, middleware := range m.middlewares {
		handler = middleware(handler)
	}
	m.Handle(pattern, handler)
}

func (m *CustomMux) HandleRouteFunc(pattern string, handler http.HandlerFunc) {
	m.HandleRoute(pattern, handler)
}

func (m *CustomMux) Use(middlewares ...func(http.Handler) http.Handler) *CustomMux {
	m.middlewares = append(m.middlewares, middlewares...)
	return m
}

func (m *CustomMux) AddGroup(pattern string, setup func(ng *Group)) {
	newGroup := &Group{
		CustomMux: &CustomMux{
			ServeMux: http.NewServeMux(),
		},
		pattern: pattern,
	}

	newGroup.action = func(ng *Group) http.Handler {
		setup(ng)

		// /api/group/ should become /api/group therefor stripping before returning
		// this way we will get "/api/group/endpoint" instead of "/api/group//endpoint"
		return http.StripPrefix(pattern[:len(pattern)-1], ng)
	}

	m.groups = append(m.groups, newGroup)
}

func CreateAndSetup(setup func(this *CustomMux) *CustomMux) *CustomMux {
	mux := &CustomMux{
		ServeMux: http.NewServeMux(),
	}

	mux = setup(mux)

	for _, g := range mux.groups {
		h := g.action(g)
		mux.HandleRoute(g.pattern, h)
	}

	return mux
}
