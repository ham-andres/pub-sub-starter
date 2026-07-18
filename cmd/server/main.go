package main

import (
	"fmt"
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
	"os/signal"
	"os"
)

func main() {
	fmt.Println("Starting Peril server...")
	const connectString = "amqp://guest:guest@localhost:5672/"
	connectStatus, err := amqp.Dial(connectString)
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer connectStatus.Close()
	fmt.Println("Connection Succesfull: localhost:5672")

	// wait for ctrl + C 
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("program shutting down and Closing the Connection")
	return	

}
