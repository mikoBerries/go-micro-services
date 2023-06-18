package main

import (
	"net/http"
)

// Broker function for routing to other microservices
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payloadBody := jsonResponse{
		Error:   false,
		Message: "this is broker route",
		// Data: nil,
	}
	app.writeJSON(w, http.StatusOK, payloadBody)
	// out, _ := json.MarshalIndent(payloadBody, "", "\t")
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusAccepted)
	// w.Write(out)
}
