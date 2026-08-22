package main

import (
	"fmt"
	"log"
	"os/signal"
	"os"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	amqp "github.com/rabbitmq/amqp091-go"

)

func main() {
	fmt.Println("Starting Peril server...")

	const connectString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connectString)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()
	fmt.Println("Server Connection Succesfull: localhost:5672")
	//connection channel
	connChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("Connection Channel creation failed: %v", err)
	}
	err = pubsub.PublishJSON(connChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{ IsPaused:	true,})
	if err != nil {
		log.Fatalf("Couldnt publish json: %v", err)
	}

	// decoupling
	gamelogic.PrintServerHelp()
	for {
		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue
		} else {
			switch inputs[0] {
			case "pause":
		    log.Println("You have Paused")
    		err = pubsub.PublishJSON(connChan,
        	routing.ExchangePerilDirect,
        	routing.PauseKey,
        	routing.PlayingState{IsPaused: true},
    		)
    		if err != nil {
        	log.Printf("could not publish pause message: %v", err)
    		}
			case "resume":
		    log.Println("Resumed")
    		err = pubsub.PublishJSON(connChan,
        	routing.ExchangePerilDirect,
	        routing.PauseKey,
  	      routing.PlayingState{IsPaused: false},
    		)
		    if err != nil {
    	    log.Printf("could not publish resume message: %v", err)
		    }
				case "quit":
			    log.Println("goodbye")
			    return
				default:
			    log.Println("Invalid command! Use Proper Command!!")
				}	
			}
	}

	// wait for ctrl + C 
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("program shutting down and Closing the Connection")
	return	

}
