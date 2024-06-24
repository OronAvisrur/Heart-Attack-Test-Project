package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// routes sets up the application's routes and middleware
func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// Specify who is allowed to connect using CORS
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},                                   // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},                 // Allow specified HTTP methods
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // Allow specified headers
		ExposedHeaders:   []string{"Link"},                                                    // Expose specified headers
		AllowCredentials: true,                                                                // Allow credentials
		MaxAge:           300,                                                                 // Max age for preflight requests
	}))

	// Allow use in the future to check if the service is still alive
	// Built-in function in go-chi
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/send", app.SendMail)

	return mux
}
