// Package repository provides functions to interact with shop logs collection in shop logs module
package repository

import (
	"context"

	constants "github.com/AlladinDev/Shopy/Constants"
	interfaces "github.com/AlladinDev/Shopy/ShopLogsModule/Interfaces"
	models "github.com/AlladinDev/Shopy/ShopLogsModule/Models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	MongoDBDatabase *mongo.Database
}

// now enforce compile time safety so that this interface implements all methods of IShopRepository
var _ interfaces.IRepository = (*Repository)(nil)

func ReturnNewRepository(mongoDBDatabase *mongo.Database) *Repository {
	return &Repository{
		MongoDBDatabase: mongoDBDatabase,
	}
}

func (repo *Repository) AddLog(ctx context.Context, logDetails models.ShopRegistrationLogs) (*mongo.InsertOneResult, error) {
	return repo.MongoDBDatabase.Collection(constants.ShopLogsCollection).InsertOne(ctx, logDetails)
}

func (repo *Repository) UpdateLog(ctx context.Context, shopID primitive.ObjectID, logDetails models.ShopRegistrationLogs) *mongo.SingleResult {
	return repo.MongoDBDatabase.Collection(constants.ShopLogsCollection).FindOneAndUpdate(ctx, bson.M{"shopId": shopID}, bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "progress", Value: logDetails.Progress[0]},
		}},
	})
}

func (repo *Repository) GetAllLogs(ctx context.Context) ([]models.ShopRegistrationLogs, error) {
	var shopLogs []models.ShopRegistrationLogs
	cur, err := repo.MongoDBDatabase.Collection(constants.ShopLogsCollection).Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &shopLogs); err != nil {
		return nil, err
	}

	return shopLogs, nil
}
