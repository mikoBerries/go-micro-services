package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"

	"github.com/MikoBerries/go-micro-services/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort = "80"
	rpcPort = "5001"
	// mongoUrl = "mongodb://mongodatabase:27017"
	gRpcPort = "5001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	log.Println("Starting sleep ")
	time.Sleep(15 * time.Second)

	log.Println("Starting logger services")
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	log.Println("Connected to mongo !")
	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()
	// closing connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panic(err)
		} else {
			log.Println("exit")
		}
	}()
	app := Config{
		Models: data.New(client),
	}
	// RPPC serve
	rpcServer := RPCServer{
		Models: data.New(client),
	}
	// register own rpc server in standart golang rpc package
	rpc.Register(rpcServer)
	go app.rpcListen()

	app.serve()

}

// rpcListen to start listening rpc
func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on :", rpcPort)

	// Start listening TCP from port
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listener.Close()

	// Start looping and handle incoming request
	for {

		rpcConn, err := listener.Accept()
		if err != nil {
			log.Println("error happend went acppting incoming RPC:", err)
			continue
		}

		go rpc.ServeConn(rpcConn)
	}

}

// serve serve and listen use http
func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	log.Printf("logger service server at port :%s \n", webPort)
	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	mongoUrl := os.Getenv("mongoUrl")

	// Create connection options
	clientOptions := options.Client().ApplyURI(mongoUrl)

	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	// Make coonection to mongo
	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
