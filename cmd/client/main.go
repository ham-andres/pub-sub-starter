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
		log.Printf("SubscribeJSON failed: %v", err)
	}

	//subcriberJSON for move
	moveQueName := "army_moves" + "." + clientName
	armyMovesKey := routing.ArmyMovesPrefix + ".*"
	err = pubsub.SubscribeJSON(conn,
														routing.ExchangePerilTopic,
														moveQueName,
														armyMovesKey,
														pubsub.Transient,
														handlerMove(gState),
													)


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
				log.Printf("Command Spawn issue: %v", err)
			}
		case "move":
			move, err := gState.CommandMove(inputs)
			if err != nil {
				log.Printf("Command move failed: %v",err)
			} else {
				routingKey := routing.ArmyMovesPrefix + "." + clientName
				err = pubsub.PublishJSON(publishCh,
																routing.ExchangePerilTopic,
																routingKey,
																move,
															)
				if err != nil {
					log.Printf("Publishing Move failed: %v", err)
				}
				log.Printf("%v Move Succesfull", move)
			}
		case "status":
			gState.CommandStatus()
		
		case "help":
			gamelogic.PrintClientHelp()
	
		case "spam":
			log.Println("Spamming not allowed yet!")
	
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			log.Println("Invalid command! Use proper command!!")
		}
	}

	// wait for ctrl + C 
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan


	

}
