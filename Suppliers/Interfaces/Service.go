package interfaces

import (
	"context"

	model "github.com/AlladinDev/Shopy/Suppliers/Model"

	"go.mongodb.org/mongo-driver/mongo"
)

type IService interface {
	RegisterSupplier(ctx context.Context, shopID string, supplierDetails model.Supplier) (*mongo.InsertOneResult, error)
	GetAllSuppliers(ctx context.Context, page int, limit int) ([]model.Supplier, error)
	GetSupplierByName(ctx context.Context, supplierName string) (model.Supplier, error)
}
