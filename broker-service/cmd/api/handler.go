package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

// HandleSubmission ahndling all routing from front end to each specified servies
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	// fetch payload
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Switch action which service will used
	switch requestPayload.Action {
	case "auth":
		// do forwarding reqeust to auth services
		log.Printf("incoming auth request : %+v \n", requestPayload)
		app.authenticate(w, requestPayload.Auth)
	default:
		app.errorJSON(w, errors.New("unkown action"))
	}

}

// authenticate forwading and recording/manipulating response from authentication-services
func (app *Config) authenticate(w http.ResponseWriter, payload AuthPayload) {

	// create json data to auth servc
	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	// Build a new Request
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// make client and execute request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// Response status checking
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credential"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("internal server error"))
		return
	}

	// Decode incoming response from auth servc
	var jsonFromService jsonResponse
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// manipulate response from auth servc
	// checking response error
	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	// write response to front
	responsePayload := jsonResponse{
		Error:   false,
		Message: "Authenticate !",
		Data:    jsonFromService.Data,
	}
	app.writeJSON(w, http.StatusAccepted, responsePayload)

}
