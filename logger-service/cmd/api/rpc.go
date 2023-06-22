package main

import (
	"context"
	"log"
	"time"

	"github.com/MikoBerries/go-micro-services/logger-service/data"
)

// RPCServer is the type for our RPC Server. Methods that take this as a receiver are available
// over RPC, as long as they are exported.
type RPCServer struct {
}

type RPCPayload struct {
	Name string
	Data string
}

// LogInfo for RPC server to logging log info to mongoDB
func (r *RPCServer) LogInfo(requestPayload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:     requestPayload.Name,
		Data:     requestPayload.Data,
		CreateAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	// resp is the message sent back to the RPC caller
	*resp = "Processed payload via RPC:" + requestPayload.Name

	return nil
}
