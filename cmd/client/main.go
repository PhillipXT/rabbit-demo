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

	gs := gamelogic.NewGameState(name)

	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, name)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient, handlerPause(gs))
	if err != nil {
		log.Fatalf("Error subscribing to queue: %v", err)
	}

	queueName = fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, name)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, queueName, routing.ArmyMovesPrefix+".*", pubsub.Transient, handlerMove(gs))
	if err != nil {
		log.Fatalf("Error subscribing to queue: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Could not create channel: %v", err)
	}

	for {
		cmd := gamelogic.GetInput()
		if len(cmd) == 0 {
			continue
		}
		switch cmd[0] {
		case "spawn":
			err := gs.CommandSpawn(cmd)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
		case "move":
			move, err := gs.CommandMove(cmd)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Printf("move %v 1", move.ToLocation)
			err = pubsub.PublishJSON(ch, routing.ExchangePerilTopic, routing.ArmyMovesPrefix+".*", move)
			if err != nil {
				log.Fatalf("Could not publish move: %v", err)
			}
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Unknown command.")
		}
	}
}
