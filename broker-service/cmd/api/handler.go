package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/rpc"
	"time"

	"github.com/MikoBerries/go-micro-services/broker-service/event"
	logs "github.com/MikoBerries/go-micro-services/broker-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// RequestPayload universa request embeded with Specific payload nedeed for each action
type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

// AuthPayload payload used to authentication-services
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LogPayload payload used to logger-service
type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// Broker function for routing to other microservices
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payloadBody := jsonResponse{
		Error:   false,
		Message: "this is broker route",
		// Data: nil,
	}
	app.writeJSON(w, http.StatusOK, payloadBody)
}

// HandleSubmission is the main point of entry into the broker. It accepts a JSON
// payload and performs an action based on the value of "action" in that JSON.
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	// fetch payload
	var requestPayload RequestPayload
	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	log.Printf("incoming  request : %+v \n", requestPayload)

	// Switch action which service will used
	switch requestPayload.Action {
	case "auth":
		// do forwarding reqeust to auth services
		app.authenticate(w, requestPayload.Auth)
	case "log":
		// use this to logging via json
		// app.logItem(w, requestPayload.Log)

		// logging via rabbitmq
		app.logEventViewRabbitMQ(w, requestPayload.Log)

		// logging via RPC
		// app.logItemViaRPC(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unkown action"))
	}

}

// logItem for forwarding to logger services
func (app *Config) logItem(w http.ResponseWriter, payload LogPayload) {
	// create json data to logger services
	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	// loger service url defined inside container network
	loggerServiceUrl := "http://logger-service/log"

	// Build a new Request
	request, err := http.NewRequest("POST", loggerServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	request.Header.Set("Context-Type", "application/json")

	// make client and execute request
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// manipulate response from logger service
	// checking response error
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, err)
		return
	}

	// write response to front
	responsePayload := jsonResponse{
		Error:   false,
		Message: "new data are logged !",
		// Data:    nil,
	}
	app.writeJSON(w, http.StatusAccepted, responsePayload)

}

// authenticate forwading and recording/manipulating response from authentication-services
func (app *Config) authenticate(w http.ResponseWriter, payload AuthPayload) {

	// create json data to auth service
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

// sendMail Forwading to mail-service to send e-mail
func (app *Config) sendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	// mail service url defined inside container network
	mailServiceURL := "http://mailer-service/sendMail"

	// Build a new Request
	request, err := http.NewRequest("POST", mailServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	request.Header.Set("Context-Type", "application/json")

	// make client and execute request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// Response status checking
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error when calling mailer services"))
		return
	}

	// write response to front
	responsePayload := jsonResponse{
		Error:   false,
		Message: "Massage send To :" + msg.To,
	}
	app.writeJSON(w, http.StatusAccepted, responsePayload)
}

// logEventViewRabbitMQ push logger event to rabbitMQ
func (app *Config) logEventViewRabbitMQ(w http.ResponseWriter, l LogPayload) {
	// try publishing to queue
	err := app.pushToQueue(l.Name, l.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// if no err make resp
	// write response to front
	responsePayload := jsonResponse{
		Error:   false,
		Message: "Logged via rabbitMq publish and listener queue",
	}
	app.writeJSON(w, http.StatusAccepted, responsePayload)
}

// pushToQueue publish given name and msg to queue
func (app *Config) pushToQueue(name string, msg string) error {
	// create emmiter object
	emmiter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}
	// Populate logger data to push in queue
	payload := LogPayload{
		Name: name,
		Data: msg,
	}
	// make json format
	jsonEvent, _ := json.MarshalIndent(&payload, "", "\t")
	// push to queue
	emmiter.Push(string(jsonEvent), "log.INFO")

	return nil
}

type RPCPayload struct {
	Name string
	Data string
}

// logItemViaRPC calling logger-service using RPC instead JSON
func (app *Config) logItemViaRPC(w http.ResponseWriter, l LogPayload) {
	// dialing our service insede container with port
	// Dial tcp => logger-service:5001 where rpc server are listening
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// populate payload to satified logger rpc payload
	rpcPayload := RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	// result serve as response from RPCServer from logger-service || Response can be ANY type even struct
	var result string
	// Starting call rpc server with function declared in server (function in)
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// start writing response to front-end
	responsePayload := jsonResponse{
		Error:   false,
		Message: result,
	}
	app.writeJSON(w, http.StatusAccepted, responsePayload)
}

// LogViaGRPC habdler function for loggin gRPC
func (app *Config) LogViaGRPC(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	conn, err := grpc.Dial("logger-service:50001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer conn.Close()

	c := logs.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = c.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: requestPayload.Log.Name,
			Data: requestPayload.Log.Data,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "logged"

	app.writeJSON(w, http.StatusAccepted, payload)
}
