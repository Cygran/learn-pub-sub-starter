package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](connection *amqp.Connection, exchange, queueName, key string, simpleQueueType SimpleQueueType, handler func(T)) error {
	ch, queue, err := DeclareAndBind(connection, exchange, queueName, key, simpleQueueType)
	if err != nil {
		return fmt.Errorf("unable to declare and bind: %s", err)
	}

	deliveries, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %s", err)
	}

	go func() {
		for delivery := range deliveries {
			var msg T
			err := json.Unmarshal(delivery.Body, &msg)
			if err != nil {
				fmt.Println("Unable to parse message")
				delivery.Ack(false)
				continue
			}
			handler(msg)
			delivery.Ack(false)
		}
	}()
	return nil
}
