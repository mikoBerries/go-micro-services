package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *config) Broker(w http.ResponseWriter, r *http.Request) {
	payloadBody := jsonResponse{
		Error:   false,
		Message: "this is broker route",
		// Data: nil,
	}
	out, _ := json.MarshalIndent(payloadBody, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}
