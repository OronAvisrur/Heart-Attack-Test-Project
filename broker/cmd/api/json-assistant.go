package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// readJSON tries to read the body of the http request and convert it into JSON
func (app *Config) readJSON(write http.ResponseWriter, read *http.Request, data any) error {
	maxBytes := 1048576 // one megabyte

	// Limit the size of the request body to one megabyte
	read.Body = http.MaxBytesReader(write, read.Body, int64(maxBytes))

	// Decode the JSON request body into the provided data structure
	decoded_data := json.NewDecoder(read.Body)
	possible_error := decoded_data.Decode(data)
	if possible_error != nil {
		return possible_error
	}

	// Check if there's any additional JSON data after the first JSON object
	possible_error = decoded_data.Decode(&struct{}{})
	if possible_error != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// writeJSON takes a response status code and arbitrary data and writes a JSON response to the client
func (app *Config) writeJSON(write http.ResponseWriter, status int, data any, headers ...http.Header) error {
	// Marshal the data into a JSON byte slice
	out, possible_error := json.Marshal(data)
	if possible_error != nil {
		return possible_error
	}

	// Add any additional headers provided
	if len(headers) > 0 {
		for key, value := range headers[0] {
			write.Header()[key] = value
		}
	}

	// Set the content type to JSON and write the status code
	write.Header().Set("Content-Type", "application/json")
	write.WriteHeader(status)

	// Write the JSON response to the response writer
	_, possible_error = write.Write(out)
	if possible_error != nil {
		return possible_error
	}

	return nil
}

// errorJSON takes an error and optionally a response status code and generates and sends
// a JSON error response
func (app *Config) errorJSON(write http.ResponseWriter, possible_error error, status ...int) error {
	// Default status code is Bad Request (400)
	statusCode := http.StatusBadRequest

	// If a status code is provided, use it instead
	if len(status) > 0 {
		statusCode = status[0]
	}

	// Create a JSON response with the error message
	var payload jsonResponse
	payload.Error = true
	payload.Message = possible_error.Error()

	// Write the JSON error response
	return app.writeJSON(write, statusCode, payload)
}
