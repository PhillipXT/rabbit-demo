package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/phillipxt/rabbit-demo/internal/pubsub"
	"github.com/phillipxt/rabbit-demo/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cstr := "amqp://guest:guest@localhost:5672"

	conn, err := amqp.Dial(cstr)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}

	defer conn.Close()

	fmt.Println("Game server connected to RabbitMQ!")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Could not create channel: %v", err)
	}

	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
	if err != nil {
		log.Fatalf("Could not publish message: %v", err)
	}

	// Wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("RabbitMQ connection closed..")
}
