package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct{}

const connection_port = "80"

func main() {

	app := Config{}

	// Print a message to the log indicating the service is starting
	log.Println("Starting knn service on port", connection_port)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", connection_port),
		Handler: app.routes(),
	}

	possible_error := server.ListenAndServe()
	if possible_error != nil {
		log.Panic(possible_error)
	}
}
