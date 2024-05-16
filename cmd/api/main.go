package main

import (
	system "booksapi/api/resource/system"
	"booksapi/api/router"
	"fmt"
	"net/http"
	"time"
)

func main() {
	buildTime := time.Now().String() // TO DO

	systemApi := system.New(buildTime)

	router := router.New(systemApi)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", 6012),
		Handler: router,
	}

	server.ListenAndServe()
}
