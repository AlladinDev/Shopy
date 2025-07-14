// Package worker contains workers for different tasks
package worker

import (
	constants "UserService/Constants"
	config "UserService/Pkg/Config"
	models "UserService/ShopLogsModule/Models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckFailedShopRegistrationEvents(rabbitConn *amqp091.Connection) {

	channel, err := rabbitConn.Channel()

	if err != nil {
		fmt.Print(err)
		return
	}

	amqpDelivery, err := channel.Consume(constants.ShopRegistrationQueueForUserModule, "", false, false, false, false, nil)

	if err != nil {
		fmt.Print(err)
		return
	}

	maxRetries := 6
	retryCount := 0

	for msg := range amqpDelivery {

		//if max retries exceeded stop worker it means some grave issue is in code which needs to be solved
		if retryCount > maxRetries {
			fmt.Println("Max retries reached stopping worker")
			break
		}

		//unmarshall the body
		var shopLogs []models.ShopRegistrationLogs
		if err := json.Unmarshal(msg.Body, &shopLogs); err != nil {
			//if error occurs while unmarshalling nack the message back so that it gets requed
			if err := msg.Nack(false, true); err != nil {
				fmt.Println(err)
			}

			//update retry count
			retryCount++

			//wait for 2 seconds before retrying
			time.Sleep(2 * time.Second)
			continue
		}

		UpdateDatabase(shopLogs)

		//ack the message after processing it so that it gets deleted from rabbitmq queue
		if err := msg.Ack(false); err != nil {
			fmt.Println(err)
			return
		}

	}
}

func UpdateDatabase(data []models.ShopRegistrationLogs) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var updateModels []mongo.WriteModel
	for _, event := range data {
		userMongoDBID, err := primitive.ObjectIDFromHex(string(event.ShopID.Hex()))

		if err != nil {
			return
		}

		filter := bson.M{"shop": userMongoDBID}                         // match documents with this shopId
		update := bson.M{"$set": bson.M{"shop": primitive.NilObjectID}} // set to null ObjectID

		model := mongo.NewUpdateManyModel().
			SetFilter(filter).
			SetUpdate(update)

		updateModels = append(updateModels, model)
	}

	result, err := config.MongoDbDatabase.Collection(constants.UserCollection).BulkWrite(ctx, updateModels, options.BulkWrite().SetOrdered(false))
	if err != nil {
		log.Fatal("Bulk delete error:", err)
	}

	log.Printf("Deleted %d user documents\n", result.ModifiedCount)
}

func CheckFailedShopRegistrationEventsWorker(td time.Duration, rabbitConn *amqp091.Connection) {
	fmt.Println("started worker for checking failed shop registration events in user module at time", time.Now().Local().Format("2006-01-02 15:04:05"))
	go CheckFailedShopRegistrationEvents(rabbitConn)
}
