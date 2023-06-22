package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// connecting to rabbit mq
	rabbitConnection, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConnection.Close()

	app := Config{
		Rabbit: rabbitConnection,
	}
	// make server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	log.Printf("starting broker services at :%s\n", webPort)
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// connect to get connection to rabbitMQ
func connect() (*amqp.Connection, error) {
	var counts int64
	interval := 5 * time.Second
	var connection *amqp.Connection
	for {
		// try connect to Rabbit MQ (Driver :// Username : Password @ Host)
		c, err := amqp.Dial(os.Getenv("rabbitMQDial"))

		if err != nil {
			log.Println("Rabbit MQ are not yet ready")
			counts++
		} else {
			log.Println("Connected to Rabbit MQ ")
			connection = c
			break

		}
		if counts > 10 {
			return nil, err
		}

		time.Sleep(interval)
	}

	return connection, nil
}
