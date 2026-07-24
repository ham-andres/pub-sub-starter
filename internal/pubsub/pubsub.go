package pubsub

import (
	"fmt"
	"encoding/json"
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	valData, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("Marshaling failed: %w", err)
	}

	msg := amqp.Publishing {
		ContentType:	"application/json",
		Body:	valData,
	}
	err = ch.PublishWithContext(context.Background(),exchange, key, false, false, msg)
	if err != nil {
		return fmt.Errorf("channel publish failed: %w", err)
	}
	return nil
}
