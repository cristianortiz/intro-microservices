package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

// declareExchange call de amqp ExchangeDeclare method wich returns an error type
// this way the required error type in consumer.setup() is returned
func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //name
		"topic",
		true,  //durable
		false, //auto-deleted?
		false, //internal?
		false, // no wait?
		nil,   //arguments?
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {

	return ch.QueueDeclare(
		"",    //name?
		false, //durable?
		false, // delete when not used
		true,  //exclusive
		false, //no-wait?
		nil,   //arguments?
	)
}
