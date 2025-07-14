// Package interfaces define interface for repository of shopLogs module
package interfaces

import (
	"context"

	models "github.com/AlladinDev/Shopy/ShopLogsModule/Models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRepository interface {
	AddLog(ctx context.Context, logDetails models.ShopRegistrationLogs) (*mongo.InsertOneResult, error)
	UpdateLog(ctx context.Context, shopID primitive.ObjectID, logDetails models.ShopRegistrationLogs) *mongo.SingleResult
	GetAllLogs(ctx context.Context) ([]models.ShopRegistrationLogs, error)
}
