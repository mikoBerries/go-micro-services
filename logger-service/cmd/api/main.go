package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MikoBerries/go-micro-services/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoUrl = "mongodb://mongo:27017"
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
	log.Println("1")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	log.Println("2")

	defer cancel()
	// closing connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panic(err)
		} else {
			log.Println("exit")
		}
	}()
	log.Println("3")

	app := Config{
		Models: data.New(client),
	}
	// app.serve()
	log.Println("4")

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Println("5")

	log.Printf("logger service server at port :%s \n", webPort)

	log.Println("6")
	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

// func (app *Config) serve() {
// 	srv := &http.Server{
// 		Addr:    fmt.Sprintf(":%s", webPort),
// 		Handler: app.routes(),
// 	}
// 	err := srv.ListenAndServe()
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

func connectToMongo() (*mongo.Client, error) {
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
