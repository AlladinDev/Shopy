// Package worker provides various workers for various functions
package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	constants "github.com/AlladinDev/Shopy/Constants"
	config "github.com/AlladinDev/Shopy/Pkg/Config"
	models "github.com/AlladinDev/Shopy/ShopLogsModule/Models"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckFailedShopRegistrationEvents(rabbitConn *amqp091.Connection) {
	channel, err := rabbitConn.Channel()

	//close the channel at last
	defer func() {
		if err := channel.Close(); err != nil {
			log.Println("failed to close channel", err)
		}
	}()

	if err != nil {
		fmt.Println(err)
		return
	}

	amqpDelivery, chanErr := channel.Consume(constants.ShopRegistrationQueueForShopModule, "", false, false, false, false, nil)

	if chanErr != nil {
		fmt.Println(chanErr)
		return
	}

	for msg := range amqpDelivery {
		//unmarshall the body data
		var shopRegistrationLogs []models.ShopRegistrationLogs
		if err := json.Unmarshal(msg.Body, &shopRegistrationLogs); err != nil {
			//if error while unmarshalling nack to requeue it back
			if err := msg.Nack(false, true); err != nil {
				fmt.Println(err)
			}
		}

		//now update database
		UpdateDatabase(shopRegistrationLogs, msg)

		//now here ack the message so that it gets deleted from rabbitmq also
		if err := msg.Ack(false); err != nil {
			fmt.Println("failed to ack message of failed shop registration", err)
		}

	}

}

func UpdateDatabase(data []models.ShopRegistrationLogs, msg amqp091.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	//use bulk write to delete multiple matching documents because data is an array
	var bulkWriteModel []mongo.WriteModel
	for _, doc := range data {
		//first convert string shopid to mongodb id format then find using that id
		shopMongoID, err := primitive.ObjectIDFromHex(doc.ShopID.Hex())

		if err != nil {
			//if error nack the message back
			if err := msg.Nack(false, true); err != nil {
				fmt.Println("Failed to nack  message when error occured while converting string shopid to mongodb id format", err)
				continue
			}
		}

		deleteOneModel := mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": shopMongoID})
		bulkWriteModel = append(bulkWriteModel, deleteOneModel)
	}

	//now here do the bulk delete
	res, err := config.MongoDbDatabase.Collection(constants.ShopCollection).BulkWrite(ctx, bulkWriteModel, options.BulkWrite().SetOrdered(false))

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Documents deleted in shop module are", res.DeletedCount)
}

func CheckFailedShopRegistrationEventsWorker(td time.Duration, rabbitConn *amqp091.Connection) {
	fmt.Println("started worker for checking failed shop registration events in shop module at time", time.Now().Local().Format("2006-01-02 15:04:05"))
	go CheckFailedShopRegistrationEvents(rabbitConn)
}
