package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// Authenticate to autheticate user & password in database (postgres SQL)
func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}
	// validate hash password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}
	// Add some logger data to logger service to this user data
	err = app.logItems("authentication", fmt.Sprintf("%s are logged in", user.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// write response
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logItems(name string, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	// create json data to logger services
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// url inside container network
	loggerServiceUrl := "http://logger-services/log"

	// Build a new Request
	request, err := http.NewRequest("POST", loggerServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	request.Header.Set("Context-Type", "application/json")

	// make client and execute request
	client := &http.Client{}

	_, err = client.Do(request)
	if err != nil {
		return err
	}
	// defer response.Body.Close()

	return nil
}
