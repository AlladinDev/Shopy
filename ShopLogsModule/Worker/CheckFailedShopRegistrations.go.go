// Package worker contains workers for different tasks
package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	constants "github.com/AlladinDev/Shopy/Constants"
	config "github.com/AlladinDev/Shopy/Pkg/Config"
	models "github.com/AlladinDev/Shopy/ShopLogsModule/Models"
	publishers "github.com/AlladinDev/Shopy/ShopLogsModule/Publishers"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
)

func CheckFailedShopRegistrationEvents(rabbitMqConn *amqp091.Connection) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var shopLogs []models.ShopRegistrationLogs

	//filter for getting documents which has initiated status but not completed and isShopRegistrationFailureNotificationSent false means that shop registration was initiated but not completed and notification is not sent yet
	filter := bson.M{
		"progress": bson.M{
			"$elemMatch": bson.M{
				"status": "initiated",
			},
		},
		"progress.status": bson.M{
			"$ne": "completed",
		},
		"isShopRegistrationFailureNotificationSent": bson.M{"$eq": false},
	}

	cur, err := config.AppConfig.MongoDatabase.Collection(constants.ShopLogsCollection).Find(ctx, filter)

	if err != nil {
		log.Println(err)
		return
	}

	if err := cur.All(ctx, &shopLogs); err != nil {
		log.Println(err)
		return
	}

	fmt.Println("len of shop logs", len(shopLogs))

	//if no shop log event is there which has no  document in progress array with status completed return back
	if len(shopLogs) == 0 {
		return
	}

	//if some shop event is there which has status initiated but not completed publish that event for consistency across all services
	channel, err := rabbitMqConn.Channel()
	if err != nil {
		return
	}

	defer func() {
		if err := channel.Close(); err != nil {
			fmt.Println("Failed to close channel created", err)
		}
	}()

	if err := publishers.PubblishToRabbitMq(channel, constants.RabbitMqShopRegistrationExchange, "", shopLogs); err != nil {
		fmt.Println("error while publishing failed shop registration events to rabbitmq", err)
		return
	}

	//now here update the status after notification is sent change its  isShopRegistrationFailureNotificationSent to true indicating that its event has been published
	//so that in next spin of this worker these shop logs wont be sent as events again
	updateRes, err := config.MongoDbDatabase.Collection(constants.ShopLogsCollection).UpdateMany(ctx, filter, bson.M{"$set": bson.M{"isShopRegistrationFailureNotificationSent": true}})
	if err != nil {
		fmt.Println("Shop registration notification emitted but failed to update shop logs", err)
		return
	}

	if updateRes.ModifiedCount == 0 {
		fmt.Println("Shop registration notification emitted but failed to update shop logs modified count is 0")
		return
	}

	fmt.Println("Successfully published failed shop events to rabbitmq and updated shop logs also")
}

func ShopRegistrationFailureCheckingWorker(td time.Duration, rabbitMqConn *amqp091.Connection) {
	ticker := time.NewTicker(td)

	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("started worker for checking failed shop registration events in shopLogs module at time", time.Now().Local().Format("2006-01-02 15:04:05"))

		go CheckFailedShopRegistrationEvents(rabbitMqConn)
	}
}
