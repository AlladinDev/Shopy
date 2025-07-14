// Package repository provides repository functions to interact with Supplier Database
package repository

import (
	constants "UserService/Constants"
	interfaces "UserService/Suppliers/Interfaces"
	model "UserService/Suppliers/Model"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	DB *mongo.Database
}

// this will ensure that this Repository struct implements all methods of IRepository interface
var _ interfaces.IRepository = (*Repository)(nil)

func ReturnNewRepository(mongodb *mongo.Database) *Repository {
	return &Repository{
		DB: mongodb,
	}
}

// AddSupplier this function register supplier
func (repo *Repository) AddSupplier(ctx context.Context, supplierData model.Supplier) (*mongo.InsertOneResult, error) {
	return repo.DB.Collection(constants.SupplierModel).InsertOne(ctx, supplierData)
}

func (repo *Repository) GetAllSuppliers(ctx context.Context, page int, limit int) ([]model.Supplier, error) {
	var suppliers []model.Supplier

	//make pagination option
	skip := (page - 1) * limit

	//this is the option to enable pagination
	options := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(limit))

	cur, err := repo.DB.Collection(constants.SupplierModel).Find(ctx, bson.D{}, options)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &suppliers); err != nil {
		return nil, err
	}

	return suppliers, nil
}

func (repo *Repository) GetSupplierByID(ctx context.Context, supplierID primitive.ObjectID) (model.Supplier, error) {
	var supplier model.Supplier

	if err := repo.DB.Collection(constants.SupplierModel).
		FindOne(ctx, bson.M{"_id": supplierID}).
		Decode(&supplier); err != nil {
		return model.Supplier{}, nil
	}

	return supplier, nil
}

func (repo *Repository) GetSupplierByName(ctx context.Context, supplierName string) (model.Supplier, error) {
	var supplier model.Supplier

	if err := repo.DB.Collection(constants.SupplierModel).
		FindOne(ctx, bson.M{"name": supplierName}).
		Decode(&supplier); err != nil {
		return model.Supplier{}, nil
	}

	return supplier, nil
}
