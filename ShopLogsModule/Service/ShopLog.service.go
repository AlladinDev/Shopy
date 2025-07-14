// Package service provides functions for service layer of shoplogs module
package service

import (
	"context"
	"fmt"

	interfaces "github.com/AlladinDev/Shopy/ShopLogsModule/Interfaces"
	models "github.com/AlladinDev/Shopy/ShopLogsModule/Models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Repository interfaces.IRepository
}

func ReturnNewService(repository interfaces.IRepository) *Service {
	return &Service{
		Repository: repository,
	}
}

// compile time check for
var _ interfaces.IService = (*Service)(nil)

func (sv *Service) AddLog(ctx context.Context, logDetails models.ShopRegistrationLogs) (*mongo.InsertOneResult, error) {
	return sv.Repository.AddLog(ctx, logDetails)
}

// UpdateLog func is for updating status of shop registration logs
func (sv *Service) UpdateLog(ctx context.Context, logDetails models.ShopRegistrationLogs) (*mongo.SingleResult, error) {
	//convert shop id string into mongo id
	shopMongoID, idErr := primitive.ObjectIDFromHex(string(logDetails.ShopID.Hex()))

	if idErr != nil {
		return nil, idErr
	}
	fmt.Println("shopid in shoplogs ", shopMongoID)

	mongoRes := sv.Repository.UpdateLog(ctx, shopMongoID, logDetails)

	return mongoRes, nil

}

func (sv *Service) GetAllShopLogs(ctx context.Context) ([]models.ShopRegistrationLogs, error) {
	return sv.Repository.GetAllLogs(ctx)
}
