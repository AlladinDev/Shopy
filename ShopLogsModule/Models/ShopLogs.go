// Package models provides schemas related to shop registration
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShopRegistrationProgress struct {
	EventTime   time.Time `json:"eventTime" bson:"eventTime"`
	Status      string    `json:"status" bson:"status"`
	HandlerFunc string    `json:"handlerFunc" bson:"handlerFunc"`
}

// ShopRegistrationLogs  this struct is for registering
type ShopRegistrationLogs struct {
	RegistrationDate time.Time          `json:"registrationDate" bson:"registrationDate"`
	UserID           primitive.ObjectID `json:"userId" bson:"userId"`
	ShopID           primitive.ObjectID `json:"shopId" bson:"shopId"`
	//expiry time means before this time progress should be updated with status completed otherwise workers will pick this log and emit its event for deletion
	ExpiryTime                                time.Time                  `json:"expiryTime" bson:"expiryTime"`
	IsShopRegistrationFailureNotificationSent bool                       `json:"isShopRegistrationFailureNotificationSent" bson:"isShopRegistrationFailureNotificationSent"`
	Progress                                  []ShopRegistrationProgress `json:"progress" bson:"progress"`
}
