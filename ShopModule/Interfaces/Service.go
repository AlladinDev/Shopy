// Package interfaces provides a set of functions in  service layer of shop module
package interfaces

import (
	models "UserService/ShopModule/Models"
	model "UserService/Suppliers/Model"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type IShopService interface {
	GetAllShops(ctx context.Context) ([]models.Shop, error)
	GetShopByName(ctx context.Context, shopName string) (models.Shop, error)
	RegisterShop(ctx context.Context, shopDetails models.Shop, userID string) (*mongo.InsertOneResult, error)
	AddSupplierToShop(ctx context.Context, supplierDetails model.SupplierRegistrationLogs) (*mongo.UpdateResult, error)
}
