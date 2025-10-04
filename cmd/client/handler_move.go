package main

import (
	"fmt"

	"github.com/phillipxt/rabbit-demo/internal/gamelogic"
	"github.com/phillipxt/rabbit-demo/internal/pubsub"
)

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(mv gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(mv)
		switch outcome {
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.NackDiscard
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			return pubsub.Ack
		}
		fmt.Println("error: unknown outcome")
		return pubsub.NackDiscard
	}
}
