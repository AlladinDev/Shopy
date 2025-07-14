// Package schemas provides schemas related to shop registration
package schemas

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShopRegistrationProgress struct {
	EventTime   time.Time `json:"eventTime" bson:"eventTime"`
	Status      string    `json:"status" bson:"status"`
	HandlerFunc string    `json:"handlerFunc" bson:"handlerFunc"`
}

type ShopRegistration struct {
	RegistrationDate time.Time                  `json:"registrationDate" bson:"registrationDate"`
	UserID           primitive.ObjectID         `json:"userId" bson:"userId"`
	ShopID           primitive.ObjectID         `json:"shopId" bson:"shopId"`
	Progress         []ShopRegistrationProgress `json:"progress" bson:"progress"`
}
