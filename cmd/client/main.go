package main

import (
	"fmt"
	"log"
	"os/signal"
	"os"
	
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
  "github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

func main() {
	fmt.Println("Starting Peril client...")
	const connectString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectString)
	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	publishCh, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel creation from Connection failed: %v", err)
	}
	defer conn.Close()
	fmt.Println("Client Connection Succesfull")

	clientName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("retrieving client name failed: %v", err)
	}
	
	// ##declare and bind Transient queue
	// ## we are commenting it out after using subscribeJSON as it calls DeclareAndBind internally.
	// it would be redundant to use two DeclareAndBind
	
	// queName := routing.PauseKey + "." + clientName
	// _, _, err = pubsub.DeclareAndBind(conn,
	// 							routing.ExchangePerilDirect,
	// 							queName,
	// 							routing.PauseKey,
	// 							pubsub.Transient,
	// 						)	
	// if err != nil {
	// 	log.Fatalf("Declare and bind failure: %v", err)
	// }

	// done, for commands
	gState := gamelogic.NewGameState(clientName)

	// subscribeJSON for Pause
	queName := routing.PauseKey + "." + clientName
	err =	 pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect,
															queName,
															routing.PauseKey,
															pubsub.Transient,
															handlerPause(gState),
														)
	if err != nil {
		log.Printf("Pause SubscribeJSON failed: %v", err)
	}

	//subcriberJSON for move
	moveQueName := "army_moves" + "." + gState.GetUsername()
	armyMovesKey := routing.ArmyMovesPrefix + ".*"
	err = pubsub.SubscribeJSON(conn,
														routing.ExchangePerilTopic,
														moveQueName,
														armyMovesKey,
														pubsub.Transient,
														handlerMove(gState, publishCh),
													)
	if err != nil {
		log.Printf("Move SubscribeJSON failed: %v ", err)
	}
	// subscriberJSON for war
	err = pubsub.SubscribeJSON(conn,
														routing.ExchangePerilTopic,
														routing.WarRecognitionsPrefix,
														routing.WarRecognitionsPrefix + ".*",
														pubsub.Durable,
														handlerWar(gState),
													)
	if err != nil {
		log.Printf("War SubscribeJSON failed: %v", err)
	}

	//REPL (read eval print loop)
	for {
		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue 
		}
		switch inputs[0] {
		case "spawn":
			err = gState.CommandSpawn(inputs)
			if err != nil {
				fmt.Printf("Command Spawn issue: %v", err)
				continue
			}
		case "move":
			move, err := gState.CommandMove(inputs)
			if err != nil {
				fmt.Printf("Command move failed: %v",err)
				continue
			}
			err = pubsub.PublishJSON(publishCh,
															routing.ExchangePerilTopic,
															routing.ArmyMovesPrefix +"."+move.Player.Username,
															move,	
														)
			if err != nil {
				fmt.Printf("error: %s \n",err)
			}
			fmt.Printf("moved %v unit to %s\n (Move Succesfull)",len(move.Units), move.ToLocation)
		case "status":
			gState.CommandStatus()
		
		case "help":
			gamelogic.PrintClientHelp()
	
		case "spam":
			fmt.Println("Spamming not allowed yet!")
	
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Invalid command! Use proper command!!")
		}
	}

	// wait for ctrl + C 
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan


	

}
