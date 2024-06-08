package main

import (
	"booksapi/api/database"
	"booksapi/api/resource/books"
	"booksapi/api/resource/system"
	"booksapi/api/router"
	"booksapi/config"
	"booksapi/logger"

	"booksapi/api/router/middlewares"
	"booksapi/docs"
	"fmt"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

var compileDate string

func main() {
	config.Init()
	logger.Init()
	database.Init()

	logger.Info("APPLICATION HAS STARTED")

	docs.SwaggerInfo.Title = "Books store API"
	docs.SwaggerInfo.Description = "This is a simple CRUD api implementation for educatinal purposes"
	docs.SwaggerInfo.Version = "1.0"

	conf := config.GetAppsettings().Config

	router := router.CreateAndSetup(func(this *router.CustomMux) *router.CustomMux {
		this.Use(middlewares.ContentTypeJSON)

		this.AddGroup("/api/system/", func(ng *router.Group) {
			systemApi := system.New(compileDate)

			ng.HandleRouteFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
				systemApi.HandleHealth(w, r)
			})
			ng.HandleRouteFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
				systemApi.HandleAbout(w, r)
			})
		})

		this.AddGroup("/api/", func(ng *router.Group) {
			ng.Use(middlewares.RequestID)
			ng.Use(middlewares.LogRequestResponse)

			booksApi := books.New()

			ng.HandleRouteFunc("GET /books", func(w http.ResponseWriter, r *http.Request) {
				booksApi.GetBooks(w, r)
			})

			ng.HandleRouteFunc("GET /books/{id}", func(w http.ResponseWriter, r *http.Request) {
				booksApi.GetBook(w, r)
			})

			ng.HandleRouteFunc("POST /books", func(w http.ResponseWriter, r *http.Request) {
				booksApi.AddBook(w, r)
			})

			ng.HandleRouteFunc("DELETE /books/{id}", func(w http.ResponseWriter, r *http.Request) {
				booksApi.RemoveBook(w, r)
			})

			ng.HandleRouteFunc("PATCH /books/{id}", func(w http.ResponseWriter, r *http.Request) {
				booksApi.UpdateBook(w, r)
			})

		})

		this.HandleFunc("GET /swagger/*", httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", conf.Port)),
		))

		return this

	})

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(conf.ReadTimeout * int(time.Second)),
		WriteTimeout: time.Duration(conf.WriteTimeout * int(time.Second)),
	}

	server.ListenAndServe()
}
