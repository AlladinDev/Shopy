// Package model provides model for this supplier module
package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Supplier struct {
	ID          primitive.ObjectID   `json:"_id" bson:"_id,omitempty"`
	ShopID      primitive.ObjectID   `json:"shopId" bson:"shopId"`
	Name        string               `json:"name" bson:"name" validate:"required,min=3,max=30"`
	Address     string               `json:"address" bson:"address" validate:"required,min=3,max=30"`
	PhoneNumber int                  `json:"phoneNumber" bson:"phoneNumber" validate:"required,number,min=1000000000,max=9999999999"`
	CreatedAt   time.Time            `json:"createdAt" bson:"createdAt"`
	Products    []primitive.ObjectID `json:"products" bson:"products"`
}
