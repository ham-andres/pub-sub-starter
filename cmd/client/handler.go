package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
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
func handlerMove(gs *gamelogic.GameState, publishChan *amqp.Channel) func(gamelogic.ArmyMove) pubsub.Acktype {
	return func(move gamelogic.ArmyMove) pubsub.Acktype {
		defer fmt.Print("> ")
		mOutcome := gs.HandleMove(move)
		fmt.Println("outcome:", mOutcome)
		switch mOutcome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeSamePlayer:
			return pubsub.Ack
	
		case gamelogic.MoveOutcomeMakeWar:
			err := pubsub.PublishJSON(publishChan,
															routing.ExchangePerilTopic,
															routing.WarRecognitionsPrefix + "." + gs.GetUsername(),
															gamelogic.RecognitionOfWar{
																Attacker: move.Player,
																Defender: gs.GetPlayerSnap(),
															},
														)
			if err != nil {
				fmt.Printf("Publishing Move failed: %s", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		}
		fmt.Println("error: unknown move outcome")
		return pubsub.NackDiscard
	}
}

// handler function for War 
func handlerWar(gs *gamelogic.GameState) func(war gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(war gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Printf("> ")
		warOutcome, _, _ := gs.HandleWar(war) 
		switch warOutcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			return pubsub.Ack

		}
		fmt.Println("error: unknown move outcome")
		return pubsub.NackDiscard
	}
}
