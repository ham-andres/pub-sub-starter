package main

import (
	"fmt"
	"log"
	"os/signal"
	"os"

	"github.com/hamandres/pub-sub-starter/internal/routing"
	"github.com/hamandres/pub-sub-starter/internal/pubsub"
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
	fmt.Println("Connection Succesfull: localhost:5672")
	//connection channel
	connChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("Connection Channel creation failed: %v", err)
	}
	err = pubsub.PublishJSON(connChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{ IsPaused:	true,})
	if err != nil {
		log.Fatalf("Couldnt publish json: %v", err)
	}

	// wait for ctrl + C 
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("program shutting down and Closing the Connection")
	return	

}
