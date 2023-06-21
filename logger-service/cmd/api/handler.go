package main

import (
	"net/http"

	"github.com/MikoBerries/go-micro-services/logger-service/data"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

// WriteLog Writing new log data to db (Mongo db collection (logs))
func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	// read incoming json
	var requestPayload JSONPayload
	_ = app.readJSON(w, r, &requestPayload)

	newLog := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(newLog)
	if err != nil {
		app.errorJSON(w, err)
	}

	resp := jsonResponse{
		Error:   false,
		Message: "Data logged",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}
