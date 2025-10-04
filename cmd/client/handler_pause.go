package main

import (
	"fmt"

	"github.com/phillipxt/rabbit-demo/internal/gamelogic"
	"github.com/phillipxt/rabbit-demo/internal/pubsub"
	"github.com/phillipxt/rabbit-demo/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}
