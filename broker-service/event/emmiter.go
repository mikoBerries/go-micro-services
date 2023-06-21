package event

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Emitter populate to pushing event to queue
type Emitter struct {
	connection *amqp.Connection
}

// setup Emmiter.setup() setup amqp.connection to channel
func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel)
}

// Push Emmiter.Push() publishing new Event to queue
func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	log.Println("pushing to Channel")
	// publishing a event

	err = channel.Publish(
		"logs_topic", // queue name
		severity,     // routing key
		false,
		false,
		amqp.Publishing{ // message to publish just like json
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	if err != nil {
		return err
	}
	return nil
}

// NewEventEmitter creating Emitter with given connection and setup()
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emtiter := Emitter{
		connection: conn,
	}
	// set-up emitter channel
	err := emtiter.setup()
	if err != nil {
		return Emitter{}, err
	}
	return emtiter, nil
}
