// Package worker provides various workers for various background tasks
package worker

import (
	constants "UserService/Constants"
	config "UserService/Pkg/Config"
	publishers "UserService/ShopLogsModule/Publishers"
	model "UserService/Suppliers/Model"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DeleteFailedSupplierRegistrations(data []model.SupplierRegistrationLogs) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var BulkModel []mongo.WriteModel
	for _, doc := range data {
		supplierMongoDBID, err := primitive.ObjectIDFromHex(doc.SupplierID.Hex())
		if err != nil {
			return
		}
		deleteModel := mongo.NewDeleteOneModel().SetFilter(bson.M{"_id": supplierMongoDBID})
		BulkModel = append(BulkModel, deleteModel)
	}
	res, err := config.MongoDbDatabase.Collection(constants.SupplierModel).BulkWrite(ctx, BulkModel, options.BulkWrite().SetOrdered(false))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Successfully deleted %c failed supplier registrations", res.ModifiedCount)
}

func CheckFailedSupplierRegistrationLogs() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//find query to check for supplier logs where registration was initiated but not completed and also expiry time for registration is reached
	findQuery := bson.M{
		"progress":                              bson.M{"$not": bson.M{"$elemMatch": bson.M{"status": "completed"}}},
		"expiryTime":                            bson.M{"$lte": time.Now()},
		"isRegistrationFailureNotificationSent": false,
	}

	//make retry with exponential backoff
	maxRetries := 3
	retryCount := 0
	delay := 1 * time.Second
	var retryErr error

	for retryCount <= maxRetries {

		//if any error it means retry but with some delay as defined by this delay variable
		if retryErr != nil {
			time.Sleep(delay)
			delay = delay * 2 //increase delay with each failure
		}

		var SupplierLogs []model.SupplierRegistrationLogs
		cur, err := config.MongoDbDatabase.Collection(constants.SupplierRegistrationLogsModel).Find(ctx, findQuery)
		if err != nil {
			retryCount++
			retryErr = err
			continue
		}

		if err := cur.All(ctx, &SupplierLogs); err != nil {
			retryCount++
			retryErr = err
			continue
		}

		//if supplierLogs is zero return back
		if len(SupplierLogs) == 0 {
			fmt.Println("Supplier logs for failed registrations is 0 returning back")
			return
		}

		///now as we got the supplier logs publish these logs using rabbitmq
		channel, err := config.AppConfig.RabbitMqConnection.Channel()
		if err != nil {
			retryCount++
			retryErr = err
			continue
		}

		if err := publishers.PubblishToRabbitMq(channel, constants.RabbiqMQSupplierExchange, "", SupplierLogs); err != nil {
			log.Println("Failed to publish failed supplier registrations to rabbitmq", err)
			continue
		}

		//now when message is published update supplier logs with isRegistrationFailureNotificationSent to true so that next time that message wont be picked up by this worker
		updateRes, err := config.MongoDbDatabase.Collection(constants.SupplierRegistrationLogsModel).UpdateMany(ctx, findQuery, bson.M{"$set": bson.M{"isRegistrationFailureNotificationSent": true}})
		if err != nil {
			log.Println("Supplier failed registrations published to rabbit but failed to update logs", err)
			continue
		}

		if updateRes.ModifiedCount == 0 {
			log.Println("Supplier failed registrations published to rabbit but failed to update logs as modified count is 0")
			continue
		}

		//here try to delete these suppliers from supplier collection
		DeleteFailedSupplierRegistrations(SupplierLogs)

		log.Println("Successfully published failed supplier registrations to rabbitmq and also updated logs")
		break //break the loop as task is complete both publishing and event emitting
	}

	if retryCount > maxRetries {
		log.Println("Max retries reached for check failed supplier registrations final error was", retryErr)
	}

}

// WCheckFailedSupplierRegistrations worker checks if any supplier registrations failed and will publish failure notification for ensuring consistency across database
func WCheckFailedSupplierRegistrations(td time.Duration) {
	ticker := time.NewTicker(td)
	defer ticker.Stop()

	for range ticker.C {
		CheckFailedSupplierRegistrationLogs()
	}
}
