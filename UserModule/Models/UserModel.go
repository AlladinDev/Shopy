// Package models provides model for user
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"Id" bson:"_id,omitempty"`
	Shop        primitive.ObjectID `json:"shop" bson:"shop"`
	Name        string             `json:"name" bson:"name"`
	Age         int                `json:"age" bson:"age" validate:"required,min=18,max=80"`
	Address     string             `json:"address" bson:"address" validate:"required,max=30,min=3"`
	PhoneNumber int                `json:"phoneNumber" bson:"phoneNumber" validate:"number"`
	Password    string             `json:"password" bson:"password" validate:"required,min=8,max=35"`
	Email       string             `json:"email" bson:"email" validate:"required,email"`
	UserType    string             `json:"userType" bson:"userType" validate:"required,oneof=user"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}
