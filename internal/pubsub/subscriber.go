package pubsub
// libraries
import (
	"fmt"
	"log"
	"bytes"
	"encoding/gob"

	"encoding/json"
	

	amqp "github.com/rabbitmq/amqp091-go"
)

// when declaring struct with lowercase mean 
// private class discoverable within same package ex: pubsub
type Acktype int 

const (
	Ack Acktype = iota 
	NackRequeue
	NackDiscard
)

func SubscribeJSON[T any](
    conn *amqp.Connection,
    exchange,
    queueName,
    key string,
    queueType SimpleQueueType, // an enum to represent "durable" or "transient"
    handler func(T) Acktype,
) error {
	subChan,_, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("Declare and bind failed in subscriber: %w", err)
	}

	//prefetch
	err = subChan.Qos(10, 0, false)
	if err != nil {
		return fmt.Errorf("Qos of JSON failed: %v", err)
	}


	messageDelivery, err := subChan.Consume(queueName, "", false, false, false, false, nil )
	if err != nil {
		return fmt.Errorf("channel.Consume failed: %w", err)
	}
	go func() {
		for delivery := range messageDelivery {
			var messages T
			err = json.Unmarshal(delivery.Body, &messages)
			if err != nil {
				log.Printf("failed to unmarshal messages: %v", err)
				delivery.Nack(false, false)
				continue
			}
			ackType := handler(messages)
			
			switch ackType {
			case Ack:
				delivery.Ack(false)
				log.Printf("Acktype: Ack ")
			case NackRequeue:
				delivery.Nack(false, true)
				log.Printf("Acktype: NackRequeue")
			case NackDiscard:
				delivery.Nack(false, false)
				log.Printf("Acktype: NackDiscard")
			}	
		}
	}()
	return nil
}

func SubscribeGob[T any](
		conn *amqp.Connection,
		exchange,
		queueName,
		key string,
		queueType SimpleQueueType,
		handler func(T) Acktype,
	) error {
		subChan, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
		if err != nil {
			return fmt.Errorf("Declare and Bind failed for SubscribeGob: %w", err)
		}
		// prefetch 
		err = subChan.Qos(10, 0, false)
		if err != nil {
			return fmt.Errorf("Qos Gob failed: %v", err)
		}


		// consumer 
		messageDelivery, err := subChan.Consume(queueName, "", false, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("subscribeGob channel consume failed: %w", err)
		}

		go func() {
			for delivery := range messageDelivery {
				var message T 
				buffer := bytes.NewReader(delivery.Body)
				decoder := gob.NewDecoder(buffer)
				err := decoder.Decode(&message)
				if err != nil {
					log.Printf("failed in Gob decoder: %v", err)
					delivery.Nack(false, false)
					continue
				}
				ackType := handler(message)

				switch ackType {
				case Ack:
					delivery.Ack(false)
					log.Printf("Acktype: Ack")

				case NackRequeue:
					delivery.Nack(false, true)
					log.Printf("Acktype: NackRequeue")

				case NackDiscard:
					delivery.Nack(false, false)
					log.Printf("Acktype: NackDiscard")
				}
			}
		}()
		return nil
}
