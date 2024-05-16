package main

import (
	"booksapi/api/resource/system"
	"booksapi/api/router"

	"booksapi/api/router/middlewares"
	"booksapi/docs"
	"fmt"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {

	docs.SwaggerInfo.Title = "Books store API"
	docs.SwaggerInfo.Description = "This is a simple CRUD api implementation for educatinal purposes"
	docs.SwaggerInfo.Version = "1.0"

	buildTime := time.Now().String() // TO DO
	systemApi := system.New(buildTime)

	router := router.CreateAndSetup(func(this *router.CustomMux) *router.CustomMux {
		this.Use(middlewares.ContentTypeJSON)

		this.AddGroup("/api/system/", func(ng *router.Group) {
			ng.HandleRouteFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
				systemApi.HandleHealth(w, r)
			})
			ng.HandleRouteFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
				systemApi.HandleAbout(w, r)
			})
		})

		this.HandleFunc("GET /swagger/*", httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", 6012)),
		))

		return this

	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", 6012),
		Handler: router,
	}

	server.ListenAndServe()
}
