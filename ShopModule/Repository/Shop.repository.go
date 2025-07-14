// Package repository provides different functions in shop module for interacting with db
package repository

import (
	"context"
	"errors"

	constants "github.com/AlladinDev/Shopy/Constants"
	Interfaces "github.com/AlladinDev/Shopy/ShopModule/Interfaces"
	models "github.com/AlladinDev/Shopy/ShopModule/Models"
	model "github.com/AlladinDev/Shopy/Suppliers/Model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShopRepository struct {
	MongodbDatabaseHandler *mongo.Database
}

// this makes sure all functions of Repository interface of this shop module are implemented by  this struct ShopRepository
var _ Interfaces.IShopRepository = (*ShopRepository)(nil)

// ReturnNewShopRepository func to return new ShopRepository with mongodb handler as dependency
func ReturnNewShopRepository(mongodbDatabase *mongo.Database) *ShopRepository {
	return &ShopRepository{
		MongodbDatabaseHandler: mongodbDatabase,
	}
}

// CreateShop this function will register shop in database
func (repo *ShopRepository) CreateShop(ctx context.Context, shopDetails models.Shop) (*mongo.InsertOneResult, error) {
	return repo.MongodbDatabaseHandler.Collection(constants.ShopCollection).InsertOne(ctx, shopDetails)
}

// GetAllShops this function will fetch all shops
func (repo *ShopRepository) GetAllShops(ctx context.Context) ([]models.Shop, error) {
	var shops []models.Shop

	cursor, err := repo.MongodbDatabaseHandler.
		Collection(constants.ShopCollection).
		Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	//now extract all shops using cursor and it will autoclose when all documents are retrieved and if any error return that
	if err := cursor.All(ctx, &shops); err != nil {
		return nil, err
	}

	return shops, nil
}

// GetShopByName this func will get shop by name
func (repo *ShopRepository) GetShopByName(ctx context.Context, shopName string) (models.Shop, error) {
	var shop models.Shop

	if err := repo.MongodbDatabaseHandler.Collection(constants.ShopCollection).
		FindOne(ctx, bson.M{"shopName": shopName}).
		Decode(&shop); err != nil {
		return models.Shop{}, nil
	}

	return shop, nil
}

// GetShopByUserID this func will fetch shop by userId
func (repo *ShopRepository) GetShopByUserID(ctx context.Context, userID primitive.ObjectID) (models.Shop, error) {
	var shop models.Shop

	if err := repo.MongodbDatabaseHandler.Collection(constants.ShopCollection).
		FindOne(ctx, bson.M{"owner": userID}).
		Decode(&shop); err != nil {
		return models.Shop{}, err
	}

	return shop, nil
}

// GetShopByShopID func will get shop by its id
func (repo *ShopRepository) GetShopByShopID(ctx context.Context, shopID primitive.ObjectID) (models.Shop, error) {
	var shop models.Shop

	if err := repo.MongodbDatabaseHandler.Collection(constants.ShopCollection).
		FindOne(ctx, bson.M{"_id": shopID}).
		Decode(&shop); err != nil {
		return models.Shop{}, nil
	}

	return shop, nil
}

// AddSupplier function adds supplier to shop mentioned by shopid in supplierDetails
func (repo *ShopRepository) AddSupplier(ctx context.Context, supplierDetails model.SupplierRegistrationLogs) (*mongo.UpdateResult, error) {
	//query for updating shop
	updateQuery := bson.M{
		"$push": bson.M{"suppliers": supplierDetails.SupplierID},
	}

	resp, err := repo.MongodbDatabaseHandler.Collection(constants.ShopCollection).UpdateOne(ctx, bson.M{"_id": supplierDetails.ShopID}, updateQuery)

	if err != nil {
		return nil, err
	}

	//if modified count is 0 it means no document is updated due to causes like failed to find matching document etc
	if resp.ModifiedCount == 0 {
		return nil, errors.New("modified Count is 0")
	}

	return resp, nil

}
