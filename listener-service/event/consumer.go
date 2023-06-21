package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// declareExchange set/declareexchange rabbitmq
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // Name
		"topic",      // Type
		true,         // durable
		false,        // auto delete?
		false,        // Internal
		false,        // no-wait ?
		nil,          // args
	)
}

// declareRandomQueue declaring random queue
func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"", // Name of queue

		true,  // durable
		false, // auto delete?
		false, // Exclusive ?
		false, // no-wait ?
		nil,   // args
	)
}
