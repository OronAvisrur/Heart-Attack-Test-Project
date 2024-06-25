package main

import "net/http"

type mailMessage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// This function convert the json to mailMessage and send it to SendSMTPMessage function in mailer
func (app *Config) SendMail(write http.ResponseWriter, read *http.Request) {
	var request_payload mailMessage

	// Write the json to mailMessage struct
	possible_error := app.readJSON(write, read, &request_payload)
	if possible_error != nil {
		app.errorJSON(write, possible_error)
		return
	}

	// Create new message to send
	msg := Message{
		From:    request_payload.From,
		To:      request_payload.To,
		Subject: request_payload.Subject,
		Data:    request_payload.Message,
	}

	// Send the message
	possible_error = app.Mailer.SendSMTPMessage(msg)
	if possible_error != nil {
		app.errorJSON(write, possible_error)
		return
	}

	// Return answer to the broker
	payload := jsonResponse{
		Error:   false,
		Message: "sent to " + request_payload.To,
	}

	app.writeJSON(write, http.StatusAccepted, payload)
}
