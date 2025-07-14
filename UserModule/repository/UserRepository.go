// Package repository contains various functions for interacting with database associated with user module
package repository

import (
	constants "UserService/Constants"
	contracts "UserService/Contracts"
	interfaces "UserService/UserModule/Interfaces"
	models "UserService/UserModule/Models"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	MongoDatabase *mongo.Database
}

// now enforce compile time check so that this repository implements IUserRepository methods
var _ interfaces.IUserRepository = (*UserRepository)(nil)

// NewUserRepository func will return a UserRepository struct witl MongodbDatabase as dependency
func NewUserRepository(mongoDatabase *mongo.Database) interfaces.IUserRepository {
	return &UserRepository{MongoDatabase: mongoDatabase}
}

// RegisterUser func to register user
func (repo *UserRepository) RegisterUser(ctx context.Context, userDetails models.User) (*mongo.InsertOneResult, error) {
	//now register user
	userCollection := repo.MongoDatabase.Collection(constants.UserCollection)
	return userCollection.InsertOne(ctx, userDetails)
}

// GetUserByPhoneNumber func to get user by phoneNumber
func (repo *UserRepository) GetUserByPhoneNumber(ctx context.Context, phoneNumber int) (models.User, error) {
	userCollection := repo.MongoDatabase.Collection(constants.UserCollection)

	var user models.User
	//get user and decode it into user type
	if err := userCollection.FindOne(ctx, bson.M{"email": phoneNumber}).Decode(&user); err != nil {
		return models.User{}, err
	}

	//if success return user
	return user, nil
}

// GetBulkUsers function to get all users
func (repo *UserRepository) GetBulkUsers(ctx context.Context) ([]models.User, error) {
	start := time.Now()
	cursor, err := repo.MongoDatabase.Collection(constants.UserCollection).Find(ctx, bson.D{})

	//if error while doing find operation return it
	if err != nil {
		return nil, err
	}
	fmt.Println("hi")

	//get all users by iterating over cursor
	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	fmt.Println("time taken", time.Since(start))

	//return users
	return users, nil
}

// GetUserByID func to get particular user by id
func (repo *UserRepository) GetUserByID(ctx context.Context, userID primitive.ObjectID) (models.User, error) {
	var user models.User
	if err := repo.MongoDatabase.Collection(constants.UserCollection).FindOne(ctx, userID).Decode(&user); err != nil {
		return models.User{}, nil
	}

	//return user
	return user, nil
}

// GetUserByEmail func to get user by email
func (repo *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	if err := repo.MongoDatabase.Collection(constants.UserCollection).FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return models.User{}, nil
	}

	return user, nil
}

// AddShop function is for adding shop to user document
func (repo *UserRepository) AddShop(ctx context.Context, shopDetails contracts.ShopRegistrationLogs) *mongo.SingleResult {
	return repo.MongoDatabase.Collection(constants.UserCollection).FindOneAndUpdate(ctx, bson.M{"_id": shopDetails.UserID}, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "shop", Value: shopDetails.ShopID},
		}},
	})
}
