package pubsub

import (
	"fmt"

	
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {

	connChan, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("Connection Channel failed: %w", err)
	}
	durable, autodelete, exclusive := setQueueType(queueType)
	connQueue, err := connChan.QueueDeclare(queueName,
																					durable,
																					autodelete,
																					exclusive,
																					false,
																					amqp.Table{
																						"x-dead-letter-exchange": "peril_dlx",
																					},
																				) 
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("queue declare failed: %w",err)
	}
	err = connChan.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("Queue Bind failed: %w", err)
	}

	return connChan, connQueue, nil


}
