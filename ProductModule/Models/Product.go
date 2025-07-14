// Package models provides model for product collection
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID                    primitive.ObjectID `json:"_id" bson:"_id"`
	Name                  string             `json:"name" bson:"name"`
	SKU                   string             `json:"sku" bson:"sku"`
	CreatedAt             time.Time          `json:"createdAt" bson:"createdAt"`
	Supplier              primitive.ObjectID `json:"supplier" bson:"supplier"`
	Shop                  primitive.ObjectID `json:"shop" bson:"shop"`
	UnitType              string             `json:"unitType" bson:"unitType"`
	CriticalStockQuantity int                `json:"criticalStockQuantity" bson:"criticalStockQuantity"`
	Type                  string             `json:"type" bson:"type"`
	SellingPrice          int                `json:"sellingPrice" bson:"sellingPrice"`
	CostPrice             int                `json:"costPrice" bson:"costPrice"`
	Profit                int                `json:"profit" bson:"profit"`
	Quantity              int                `json:"quantity" bson:"quantity"`
}
