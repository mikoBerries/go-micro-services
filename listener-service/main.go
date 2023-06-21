package main

import (
	"log"
	"os"
	"time"

	"github.com/MikoBerries/go-micro-services/listener-service/event"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connecting to rabbit mq
	rabbitConnection, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConnection.Close()

	// Create new consumer
	consumer, err := event.NewConsumer(rabbitConnection)
	if err != nil {
		log.Panic(err)
	}
	// declaring what topics wa want to consume
	topics := []string{"log.*"}

	// Start consumer to listening to rabbitMQ
	err = consumer.Listen(topics)
	if err != nil {
		log.Println(err)
	}
}

// connect to get connection to rabbitMQ
func connect() (*amqp.Connection, error) {
	var counts int64
	interval := 5 * time.Second
	var connection *amqp.Connection
	// dsn := os.Getenv("DSN")
	for {
		// try connect to Rabbit MQ (Driver :// Username : Password @ Host)
		// c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
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
