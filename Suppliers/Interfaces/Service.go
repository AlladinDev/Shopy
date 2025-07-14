package interfaces

import (
	model "UserService/Suppliers/Model"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type IService interface {
	RegisterSupplier(ctx context.Context, shopID string, supplierDetails model.Supplier) (*mongo.InsertOneResult, error)
	GetAllSuppliers(ctx context.Context, page int, limit int) ([]model.Supplier, error)
	GetSupplierByName(ctx context.Context, supplierName string) (model.Supplier, error)
}
