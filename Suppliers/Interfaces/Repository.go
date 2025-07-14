// Package interfaces provides interfaces for supplier module
package interfaces

import (
	model "UserService/Suppliers/Model"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IRepository interface {
	AddSupplier(ctx context.Context, supplierInfo model.Supplier) (*mongo.InsertOneResult, error)
	GetAllSuppliers(ctx context.Context, page int, limit int) ([]model.Supplier, error)
	GetSupplierByName(ctx context.Context, name string) (model.Supplier, error)
	GetSupplierByID(ctx context.Context, ID primitive.ObjectID) (model.Supplier, error)
}
