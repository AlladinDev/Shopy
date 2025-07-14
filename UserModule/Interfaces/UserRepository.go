// Package interfaces contains interfaces used by userModule
package interfaces

import (
	contracts "UserService/Contracts"
	models "UserService/UserModule/Models"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IUserRepository interface {
	RegisterUser(ctx context.Context, userDetails models.User) (*mongo.InsertOneResult, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber int) (models.User, error)
	GetBulkUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, userID primitive.ObjectID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	AddShop(ctx context.Context, shopDetails contracts.ShopRegistrationLogs) *mongo.SingleResult
}
