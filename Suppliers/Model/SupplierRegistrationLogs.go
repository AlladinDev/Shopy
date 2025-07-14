// Package model this model provides model for supplier logs which is used when registering supplier
package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//#################this model ensures that supplier info is added into shop model and shop id into supplier collection in a reliable way and if any operation fails it can revoked

type RegistrationProgress struct {
	Status    string    `json:"status" bson:"status"`
	EventTime time.Time `json:"eventTime" bson:"eventTime"`
}

type SupplierRegistrationLogs struct {
	ID               primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	SupplierID       primitive.ObjectID `json:"supplierId" bson:"supplierId"`
	ShopID           primitive.ObjectID `json:"shopId" bson:"shopId"`
	RegistrationDate time.Time          `json:"registrationDate" bson:"registrationDate"`
	//expiry time means before this expiry this log status should be changed to completed from initiated otherwise it will be considered failed event
	ExpiryTime                            time.Time              `json:"expiryTime" bson:"expiryTime"`
	IsRegistrationFailureNotificationSent bool                   `json:"isRegistrationFailureNotificationSent" bson:"isRegistrationFailureNotificationSent"`
	Progress                              []RegistrationProgress `json:"progress" bson:"progress"`
}
