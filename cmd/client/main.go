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

	name, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error getting username: %v", err)
	}

	queueName := fmt.Sprintf("pause.%s", name)
	pubsub.CreateChannel(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient)

	game := gamelogic.NewGameState(name)

loop:
	for {
		cmd := gamelogic.GetInput()
		if len(cmd) == 0 {
			continue
		}
		switch cmd[0] {
		case "spawn":
			err := game.CommandSpawn(cmd)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
		case "move":
			move, err := game.CommandMove(cmd)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Printf("move %v 1", move.ToLocation)
		case "status":
			game.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet")
		case "quit":
			gamelogic.PrintQuit()
			break loop
		default:
			fmt.Println("Unknown command.")
		}
	}

	fmt.Println("RabbitMQ connection closed..")
}
