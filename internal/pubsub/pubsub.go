package pubsub

import (
	"fmt"
	"bytes"
	"encoding/json"
	"context"
	"encoding/gob"

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

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(val)
	if err != nil {
		return fmt.Errorf("gob encoder failed: %v", err)
	}

	msg := amqp.Publishing {
		ContentType:	"application/gob",
		Body: buffer.Bytes(),
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, msg)
	if err != nil {
		return fmt.Errorf("Channel gob publish failed: %v", err)
	}
	return nil
}
