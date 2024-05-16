package main

import (
	"booksapi/api/router"
	"fmt"
	"net/http"
)

func main() {
	router := router.New()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", 6012),
		Handler: router,
	}

	server.ListenAndServe()
}
