package pubsub

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T)) error {
	ch, _, err := CreateChannel(conn, exchange, queueName, key, queueType)
	if err != nil {
		fmt.Println("Error creating channel (SubscribeJSON).")
		return err
	}

	del, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println("Error consuming channel (SubscribeJSON).")
		return err
	}

	go func() {
		for msg := range del {
			var v T
			err := json.Unmarshal(msg.Body, &v)
			if err != nil {
				fmt.Printf("Error during unmarshal: %v\n", err)
				continue
			}
			handler(v)
			if err := msg.Ack(false); err != nil {
				fmt.Printf("Error during ack: %v\n", err)
				continue
			}
		}
	}()

	return nil
}
