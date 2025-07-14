// Package contracts provides contracts ie interfaces for external api calls or contracts
package contracts

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShopRegistrationProgress struct {
	EventTime time.Time `json:"eventTime" bson:"eventTime"`
	Status    string    `json:"status" bson:"status"`

	HandlerFunc string `json:"handlerFunc" bson:"handlerFunc"`
}

// ShopRegistrationLogs  this struct is for registering
type ShopRegistrationLogs struct {
	RegistrationDate                          time.Time                  `json:"registrationDate" bson:"registrationDate"`
	UserID                                    primitive.ObjectID         `json:"userId" bson:"userId"`
	ExpiryTime                                time.Time                  `json:"expiryTime" bson:"expiryTime"`
	IsShopRegistrationFailureNotificationSent bool                       `json:"isShopRegistrationFailureNotificationSent" bson:"isShopRegistrationFailureNotificationSent"`
	ShopID                                    primitive.ObjectID         `json:"shopId" bson:"shopId"`
	Progress                                  []ShopRegistrationProgress `json:"progress" bson:"progress"`
}
