package utils

import (
	"github.com/rabbitmq/amqp091-go"
)

func ConsumeRabbitMessages(conn *amqp091.Connection, chann *amqp091.Channel, queueName string) (<-chan amqp091.Delivery, error) {

	msgs, err := chann.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
