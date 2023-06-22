package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/MikoBerries/go-micro-services/logger-service/data"
	logs "github.com/MikoBerries/go-micro-services/logger-service/pb"
	"google.golang.org/grpc"
)

// LogServer populate server for gRPC
type LogServer struct {
	logs.UnimplementedLogServiceServer // impelent all declared rpc func in proto file
	Models                             data.Models
}

// WriteLog Implementing function that writen in proto file
func (lServer *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := lServer.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// return response as address of struct logs.LogResponse
	// gRPC are similar concept as RPC
	res := &logs.LogResponse{Result: "logged!"}
	return res, nil
}

// gRPCListen to start gRPC server
func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
	// make new gRPC server
	s := grpc.NewServer()
	// Register our server that implement all function that avail to call
	logs.RegisterLogServiceServer(s, &LogServer{Models: app.Models})

	log.Printf("gRPC Server started on port %s", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
