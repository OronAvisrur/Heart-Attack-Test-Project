package main

import (
	"fmt"
	"log"
	"net/http"
)

const connection_port = "80"

type Config struct{}

func main() {

	app := Config{}

	// Print a message to the log indicating the service is starting
	log.Printf("Starting broker service on port %s \n", connection_port)

	// Define HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", connection_port), // Define the server address (port)
		Handler: app.routes(),                        // Set the request handler (Handler) to the routes method of the app instance
	}

	// Start the server and listen for requests
	possible_error := server.ListenAndServe()

	// If there is an error starting the server, log it and panic (crash the program)
	if possible_error != nil {
		log.Panic(possible_error)
	}
}
