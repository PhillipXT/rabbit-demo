package main

import (
	"fmt"
	"log"

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

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Could not create channel: %v", err)
	}

	gamelogic.PrintServerHelp()

	for {
		cmd := gamelogic.GetInput()
		if len(cmd) == 0 {
			continue
		} else if cmd[0] == "pause" {
			fmt.Println("Sending 'pause' message...")
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				log.Fatalf("Could not publish message: %v", err)
			}
		} else if cmd[0] == "resume" {
			fmt.Println("Sending 'resume' message...")
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				log.Fatalf("Could not publish message: %v", err)
			}
		} else if cmd[0] == "quit" {
			fmt.Println("Exiting game loop...")
			break
		} else {
			fmt.Println("Unknown command.")
		}
	}

	fmt.Println("RabbitMQ connection closed..")
}
