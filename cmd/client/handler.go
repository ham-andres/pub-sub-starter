package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)

// helper function for Pause subscriber
func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	} 
}

// handler function for Move subscriber
func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	return func(mo gamelogic.ArmyMove) {
		defer fmt.Print("> ")
		gs.HandleMove(mo)
	}
}
