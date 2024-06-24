package main

import (
	"errors"
	"fmt"
	"net/http"
)

type requestsPayload struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func (app *Config) Authenticate(write http.ResponseWriter, read *http.Request) {
	var requests_payload requestsPayload

	possible_error := app.readJSON(write, read, &requests_payload)

	if possible_error != nil {
		app.errorJSON(write, possible_error, http.StatusBadRequest)
		return
	}

	//validate the user against the database
	user, possible_error := app.Models.User.GetUserByName()
	if possible_error != nil {
		app.errorJSON(write, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, possible_error := user.IsPasswordMatches(requests_payload.Password)
	if possible_error != nil || !valid {
		app.errorJSON(write, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	pay_load := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.UserName),
		Data:    user,
	}

	app.writeJSON(write, http.StatusAccepted, pay_load)
}
