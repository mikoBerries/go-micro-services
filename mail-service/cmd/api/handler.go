package main

import (
	"fmt"
	"net/http"
)

// mailMassage incoming request for sendEmail
type mailMassage struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// SendEmail handling new request to sending email consume (mailMassage)
func (app *Config) SendEmail(w http.ResponseWriter, r *http.Request) {
	var requestPayload mailMassage

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Populate incoming reqeust to Message model
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}
	// execute send message use SMTP
	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// Write some response
	responsePayload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Success -  Email are sended to (%s)", msg.To),
	}
	app.writeJSON(w, http.StatusAccepted, responsePayload)
}
