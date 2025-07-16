package interfaces

import (
	"context"

	models "github.com/AlladinDev/Shopy/ShopModule/Models"
	model "github.com/AlladinDev/Shopy/Suppliers/Model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IShopRepository interface {
	GetAllShops(ctx context.Context) ([]models.Shop, error)
	CreateShop(ctx context.Context, shopDetails models.Shop) (*mongo.InsertOneResult, error)
	GetShopByName(ctx context.Context, shopName string) (models.Shop, error)
	GetShopByUserID(ctx context.Context, userID primitive.ObjectID) (models.Shop, error)
	GetShopByShopID(ctx context.Context, shopID primitive.ObjectID) (models.Shop, error)
	AddSupplier(ctx context.Context, supplierDetails model.SupplierRegistrationLogs) (*mongo.UpdateResult, error)
}
