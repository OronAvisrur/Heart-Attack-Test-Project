package main

import "net/http"

type mailMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) SendMail(write http.ResponseWriter, read *http.Request) {
	var request_payload mailMessage

	possible_error := app.readJSON(write, read, &request_payload)
	if possible_error != nil {
		app.errorJSON(write, possible_error)
		return
	}

	msg := Message{
		From:    request_payload.From,
		To:      request_payload.To,
		Subject: request_payload.Subject,
		Data:    request_payload.Message,
	}

	possible_error = app.Mailer.SendSMTPMessage(msg)
	if possible_error != nil {
		app.errorJSON(write, possible_error)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + request_payload.To,
	}

	app.writeJSON(write, http.StatusAccepted, payload)
}
