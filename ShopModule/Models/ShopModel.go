// Package models provides model for shop
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Shop struct {
	ID        primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	ShopName  string               `json:"shopName" bson:"shopName" validate:"required,min=3,max=100"`
	Address   string               `json:"address" bson:"address" validate:"required"`
	Pincode   string               `json:"pincode" bson:"pincode" validate:"required,min=5,numeric"`
	CreatedAt time.Time            `json:"createdAt" bson:"createdAt"`
	ShopType  string               `json:"shopType" bson:"shopType" validate:"required,oneof=retail wholesale"`
	Owner     primitive.ObjectID   `json:"owner" bson:"owner" validate:"required"`
	Suppliers []primitive.ObjectID `json:"suppliers" bson:"suppliers"`
}
