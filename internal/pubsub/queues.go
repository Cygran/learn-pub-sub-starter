package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	QueueTypeDurable = iota
	QueueTypeTransient
)

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, simpleQueueType int) (*amqp.Channel, amqp.Queue, error) {
	// create a channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not create channel: %s", err)
	}

	// Set queue properties based on simpleQueueType
	isDurable := simpleQueueType == QueueTypeDurable
	isTransient := simpleQueueType == QueueTypeTransient
	// Create queue
	queue, err := ch.QueueDeclare(
		queueName,
		isDurable,
		isTransient,
		isTransient,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not declare queue: %s", err)
	}
	//bind queue to exchange
	err = ch.QueueBind(
		queue.Name,
		key,
		exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not bind queue: %s", err)
	}
	return ch, queue, nil
}
