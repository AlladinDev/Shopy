package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	constants "github.com/AlladinDev/Shopy/Constants"
	config "github.com/AlladinDev/Shopy/Pkg/Config"
	model "github.com/AlladinDev/Shopy/Suppliers/Model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckFailedSupplierRegistrationLogs() {
	//use retries logic for retries
	maxRetries := 3
	delay := 1 * time.Second
	retryCount := 0
	var retryErr error

	for retryCount <= maxRetries {

		//if retryErr has some error do delay before making another attempt
		if retryErr != nil {
			time.Sleep(delay)
			delay = delay * 2
		}

		channel, err := config.RabbitConnection.Channel()
		if err != nil {
			retryCount++
			retryErr = err
			log.Println("Failed to create channel ", err)
			continue
		}

		amqpDelivery, err := channel.Consume(constants.SupplierRegistrationFailedQueue, "", false, false, false, false, nil)
		if err != nil {
			retryCount++
			retryErr = err
			log.Println("Failed to consume channel ", err)
			continue
		}

		//reinitialise retry count and delay for retrying message consumption
		retryCount = 0
		delay = 0
		fmt.Println("listening for any failed supplier registrations in shop module")
		for msg := range amqpDelivery {

			if retryErr != nil {
				time.Sleep(delay)
				delay = delay * 2
			}

			if retryCount > maxRetries {
				fmt.Println("max limit exceeded breaking from message consumtion moving messages to dlx queue", msg.Redelivered)
				if err := msg.Nack(false, false); err != nil {
					fmt.Println(err)
				}
				continue
			}

			fmt.Println("Received messages form rabbit about failed suppliers", msg.AppId)

			//first unmarshall the msg
			var supplierLogs []model.SupplierRegistrationLogs

			if err := json.Unmarshal(msg.Body, &supplierLogs); err != nil {
				retryCount++
				retryErr = err

				if err := msg.Nack(false, true); err != nil {
					log.Println("Failed to nack", err)
				}

				log.Println("error is ", err, supplierLogs)
				//wait for 1 second before consuming another msg in case when error is there

				continue
			}
			fmt.Println("message received from rabbit mq is", supplierLogs)

			//now call update db function for db operations
			err = UpdateShopDB(supplierLogs)
			if err != nil {
				retryErr = err
				retryCount++
				log.Println("UpdateShopDB failed:", err)
				if err := msg.Nack(false, true); err != nil { // ✅ Only one Nack after all retries fail
					fmt.Println(err)
				}
			} else {
				if err := msg.Ack(false); err != nil { // ✅ Only one Ack on success
					fmt.Println(err)
				}
			}
		}
	}
}

// UpdateShopDB function to update shop and remove any failed supplier registration ids
func UpdateShopDB(data []model.SupplierRegistrationLogs) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var mongoBulkModel []mongo.WriteModel
	for _, doc := range data {
		supplierID, err := primitive.ObjectIDFromHex(doc.SupplierID.Hex())
		if err != nil {
			return fmt.Errorf("invalid supplierID: %w", err)
		}

		shopMongoDBID, err := primitive.ObjectIDFromHex(doc.ShopID.Hex())
		if err != nil {
			return fmt.Errorf("invalid shopID: %w", err)
		}
		fmt.Println(supplierID, shopMongoDBID)
		updateFilter := bson.M{"$pull": bson.M{"suppliers": supplierID}}
		updateModel := mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": shopMongoDBID}).SetUpdate(updateFilter)
		mongoBulkModel = append(mongoBulkModel, updateModel)
	}

	maxRetries := 4
	delay := 1 * time.Second
	retryCount := 0
	var err error
	var bulkWriteRes *mongo.BulkWriteResult

	for retryCount <= maxRetries {
		if err != nil {
			time.Sleep(delay)
			delay *= 2
		}

		bulkWriteRes, err = config.MongoDbDatabase.Collection(constants.ShopCollection).BulkWrite(ctx, mongoBulkModel, options.BulkWrite().SetOrdered(false))
		if err == nil {
			break
		}

		log.Println("Mongo write error:", err)
		retryCount++
	}

	fmt.Println(bulkWriteRes.MatchedCount, bulkWriteRes.ModifiedCount)
	if err != nil || bulkWriteRes == nil || bulkWriteRes.ModifiedCount == 0 {
		return fmt.Errorf("mongo write failed or no documents modified: %w", err)
	}

	fmt.Println("Successfully updated shop db")
	return nil
}
