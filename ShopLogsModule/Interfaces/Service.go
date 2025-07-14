package interfaces

import (
	models "UserService/ShopLogsModule/Models"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type IService interface {
	AddLog(ctx context.Context, logDetails models.ShopRegistrationLogs) (*mongo.InsertOneResult, error)
	UpdateLog(ctx context.Context, logDetails models.ShopRegistrationLogs) (*mongo.SingleResult, error)
	GetAllShopLogs(ctx context.Context) ([]models.ShopRegistrationLogs, error)
}
