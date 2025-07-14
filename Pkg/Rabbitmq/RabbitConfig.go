// Package rabbitmq provides functions related to rabbitmq like connecting to rabbitmq
package rabbitmq

import (
	"fmt"
	"log"
	"os"

	constants "github.com/AlladinDev/Shopy/Constants"

	"github.com/rabbitmq/amqp091-go"
)

func DLXForSupplierExchange(chann *amqp091.Channel) error {
	if err := chann.ExchangeDeclare(constants.NameOfDLXForSupplier, amqp091.ExchangeDirect, true, false, false, false, nil); err != nil {
		return err
	}
	supplierDlxQueue, err := chann.QueueDeclare("DLXQueueSupplier", true, false, false, false, nil)
	if err != nil {
		return err
	}
	if err := chann.QueueBind(supplierDlxQueue.Name, "retry", constants.NameOfDLXForSupplier, false, nil); err != nil {
		return err
	}
	return nil
}

func ConnectToRabbitMq() (*amqp091.Channel, *amqp091.Connection) {
	rabbitURL := os.Getenv("RABBITMQ_URL")

	if rabbitURL == "" {
		log.Fatal("RABBITMQ_URL missing in env file")
	}

	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		log.Fatal("Error while connecting to rabbitmq", err)
	}

	//create a dummy channel for defining queues and exchanges
	channel, err := conn.Channel()

	if err != nil {
		log.Fatal("error while creating rabbitmq channel", err)
	}

	//define shop_exchange exchange
	if err := channel.ExchangeDeclare(constants.RabbitMqShopRegistrationExchange, amqp091.ExchangeFanout, true, false, false, false, nil); err != nil {
		log.Fatal("Error While Defining Shop_Exchange", err)
	}

	//now define queues

	//queue for user module to process shop registration events
	userShopRegistrationQueue, err := channel.QueueDeclare(constants.ShopRegistrationQueueForUserModule, true, false, false, false, nil)

	if err != nil {
		log.Fatal("Error While Declaring Queue for shop Registration", err)
	}

	//queue for shop module to process shop registration events
	shopRegistrationQueueShopModule, err := channel.QueueDeclare(constants.ShopRegistrationQueueForShopModule, true, false, false, false, nil)

	if err != nil {
		log.Fatal("Error While Declaring Queue for shop Registration", err)
	}

	//##################   now bind all queues with exchanges

	//now bind the shop queue with exchange
	if err := channel.QueueBind(userShopRegistrationQueue.Name, "", constants.RabbitMqShopRegistrationExchange, false, nil); err != nil {
		log.Fatal("Failed to bind shop_registration queue for user module with Shop_Exchange exchange", err)
	}

	//now bind the userShopRegistrationQueue  with exchange
	if err := channel.QueueBind(shopRegistrationQueueShopModule.Name, "", constants.RabbitMqShopRegistrationExchange, false, nil); err != nil {
		log.Fatal("Failed to bind shopRegistration  queue  for shop module with Shop_Exchange exchange", err)
	}

	//########## define exchange for supplier related events
	if err := channel.ExchangeDeclare(constants.RabbiqMQSupplierExchange, amqp091.ExchangeFanout, true, false, false, false, nil); err != nil {
		log.Fatal("Failed to declare exchange for supplier", err)
	}

	//now define queue for this supplier exchange
	supplerQueue, err := channel.QueueDeclare(constants.SupplierRegistrationFailedQueue, true, false, false, false, amqp091.Table{
		"x-dead-letter-exchange":    constants.NameOfDLXForSupplier, // Send failed messages here
		"x-dead-letter-routing-key": "retry",                        // Route key used by retry exchange
	})

	if err != nil {
		log.Fatal("Failed to declare exchange for supplier", err)
	}

	//now bind supplier queue with supplier exchange
	if err := channel.QueueBind(supplerQueue.Name, "", constants.RabbiqMQSupplierExchange, false, nil); err != nil {
		log.Fatal("failed to bind queue for supplier with supplier exchange", err)
	}

	//now configure dead letter queue for supplier exchange
	if err := DLXForSupplierExchange(channel); err != nil {
		log.Fatal(err)
	}

	fmt.Println("RabbitMq Connected  Successfully")
	//now return the channel and connection for publising consuming purpose
	return channel, conn
}
