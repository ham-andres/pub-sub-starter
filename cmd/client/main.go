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
	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queName, routing.PauseKey, pubsub.Transient)	
	if err != nil {
		log.Fatalf("Declare and bind failure: %v", err)
	}

	// wait for ctrl + C 
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan


	

}
