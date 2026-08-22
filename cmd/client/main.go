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
	defer conn.Close()
	fmt.Println("Client Connection Succesfull")

	clientName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("retrieving client name failed: %v", err)
	}
	
	queName := routing.PauseKey + "." + clientName
	_, _, err = pubsub.DeclareAndBind(conn,
								routing.ExchangePerilDirect,
								queName,
								routing.PauseKey,
								pubsub.Transient,
							)	
	if err != nil {
		log.Fatalf("Declare and bind failure: %v", err)
	}

	// current working area
	gState := gamelogic.NewGameState(clientName)
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
				log.Printf("Move Succesfull %v",move)
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
