package pubsub

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](conn *amqp.Connection, exchange, queueName, key string, queueType SimpleQueueType, handler func(T) AckType) error {
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
			switch handler(v) {
			case Ack:
				log.Println("Message acknowledged.")
				msg.Ack(false)
			case NackDiscard:
				log.Println("Message discarded.")
				msg.Nack(false, false)
			case NackRequeue:
				log.Println("Message requeued.")
				msg.Nack(false, true)
			}
		}
	}()

	return nil
}
