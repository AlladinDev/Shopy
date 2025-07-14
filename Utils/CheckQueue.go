package utils

import "github.com/rabbitmq/amqp091-go"

func CheckRabbitQueueExistence(conn *amqp091.Connection, queueName string) (*amqp091.Channel, int, error) {

	channel, err := conn.Channel()

	if err != nil {
		return nil, 0, err
	}

	queue, err := channel.QueueDeclare(queueName, true, false, false, false, nil)

	if err != nil {
		return nil, 0, err
	}

	return channel, queue.Messages, nil
}
