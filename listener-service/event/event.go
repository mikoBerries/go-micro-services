package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer as consumer to rabbit mq
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

// NewConsumer create Consumer with given conn
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

// setup setup exchange and queue to listen
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	// defer channel.Close()
	// declaring exchange
	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// decraling queue
	randomQueue, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	// Binding Queue to this channel
	for _, topic := range topics {
		err = ch.QueueBind(
			randomQueue.Name, // Queue
			topic,            // Exchange
			"logs_topic",     // RoutingKey
			false,            // NoWait
			nil,              // Arguments
		)
		if err != nil {
			return err
		}
	}
	// ch.Consume returnin <-chan of ampq.Delivery and error
	messages, err := ch.Consume(randomQueue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	log.Printf("Waiting incoming messages from [exchane ,queue] - [`logs_topic`,`%s`]\n", randomQueue.Name)
	fmt.Printf("Waiting incoming messages from [exchane ,queue] - [`logs_topic`,`%s`]\n", randomQueue.Name)
	// loop massage chan forever to listen incoming messages
	go func() {
		for message := range messages { //forr channel
			var payload Payload // Data insede of rabbitmq
			_ = json.Unmarshal(message.Body, &payload)
			// run another go rutine to faster consume so quque are not full
			go handlePayload(payload)
		}
	}()

	// Stuck this go routine using forever channel of boolean
	forever := make(chan bool)
	<-forever
	return nil
}

// handlePayload handle each payload and forwading to specifict services similar to broker
func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// for logger services
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "":

	default:

	}
}

// logEvent calling looger-service to save a new log
func logEvent(payload Payload) error {
	// create json data to logger services
	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	// url inside container network
	loggerServiceUrl := "http://logger-service/log"

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

	return nil
}
