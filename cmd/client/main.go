package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/phillipxt/rabbit-demo/internal/gamelogic"
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

	name, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error getting username: %v", err)
	}

	queueName := fmt.Sprintf("pause.%s", name)
	pubsub.CreateChannel(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient)

	// Wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("RabbitMQ connection closed..")
}
