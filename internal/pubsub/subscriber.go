package pubsub
// libraries
import (
	"fmt"
	"log"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
    conn *amqp.Connection,
    exchange,
    queueName,
    key string,
    queueType SimpleQueueType, // an enum to represent "durable" or "transient"
    handler func(T),
) error {
	subChan,_, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("Declare and bind failed in subscriber: %w", err)
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
				log.Fatalf("Unmarshaling the delivery message failed: %v", err)
			}
			handler(messages)
			delivery.Ack(false)
		}
	}()
	return nil
}
