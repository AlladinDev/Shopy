// Package publishers contains publish function for publishing events to rabbitmq
package publishers

import (
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

func PubblishToRabbitMq(rabbitChan *amqp091.Channel, exchangeName string, key string, dataToSend any) error {
	jsonData, err := json.Marshal(&dataToSend)

	if err != nil {
		return err
	}

	if err := rabbitChan.Publish(exchangeName, key, false, false, amqp091.Publishing{
		Body: jsonData,
	}); err != nil {
		return err
	}

	return nil
}
