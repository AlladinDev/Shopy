// Package mongodb provides functions to connect to mongodb
package mongodb

import (
	"context"
	"errors"
	"os"

	constants "github.com/AlladinDev/Shopy/Constants"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongodb() (*mongo.Database, error) {
	ctx, _ := context.WithCancel(context.Background())
	mongodbURL := os.Getenv("MONGODB_URL")

	if mongodbURL == "" {
		return nil, errors.New("mongodb url not found in env")
	}

	mongoOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URL"))

	mongoClient, err := mongo.Connect(ctx, mongoOptions)

	if err != nil {
		return nil, err
	}
	//make the database handler global for shared use to avoid creating handler everytime
	mongoDatabaseHandler := mongoClient.Database(constants.DatabaseName)

	return mongoDatabaseHandler, nil
}
