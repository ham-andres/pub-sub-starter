package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

// helper function for Pause subscriber
func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.Acktype {
	return func(ps routing.PlayingState) pubsub.Acktype {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	} 
}

// handler function for Move subscriber
func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(mo gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		mOutcome := gs.HandleMove(mo)
		if (mOutcome == gamelogic.MoveOutComeSafe || mOutcome == gamelogic.MoveOutcomeMakeWar) {
			return pubsub.Ack
		} else {
			return pubsub.NackDiscard
		} 
	}
}
