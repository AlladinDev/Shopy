// Package service provides service functions for supplier module
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	constants "github.com/AlladinDev/Shopy/Constants"
	config "github.com/AlladinDev/Shopy/Pkg/Config"
	interfaces "github.com/AlladinDev/Shopy/Suppliers/Interfaces"
	model "github.com/AlladinDev/Shopy/Suppliers/Model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	Repo interfaces.IRepository
}

// this line will make sure this Service struct implements all methods of IService interface
var _ interfaces.IService = (*Service)(nil)

// ReturnNewService returns new Service for using its methods
func ReturnNewService(repo interfaces.IRepository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (sv *Service) GetAllSuppliers(ctx context.Context, page int, limit int) ([]model.Supplier, error) {
	return sv.Repo.GetAllSuppliers(ctx, page, limit)
}

func (sv *Service) RegisterSupplier(ctx context.Context, shopID string, supplierDetails model.Supplier) (*mongo.InsertOneResult, error) {
	fmt.Println("1 pass", supplierDetails.ID, shopID)
	//convert shopid which is string into mongodb id format
	shopMongoID, err := primitive.ObjectIDFromHex(shopID)

	if err != nil {
		return nil, err
	}

	//generate mongoid for supplierDetails
	supplierDetails.ID = primitive.NewObjectID()

	fmt.Println("1 pass", supplierDetails.ID)

	//add shopid to this supplierdetails
	supplierDetails.ShopID = shopMongoID

	//now using outbox pattern save supplier id in shop collection and shopid in supplier collection

	//first save this log in supplierLog model for this create payload first
	supplierLogsPayload := model.SupplierRegistrationLogs{
		RegistrationDate: time.Now(),
		SupplierID:       supplierDetails.ID,
		ExpiryTime:       time.Now().Add(10 * time.Second),
		ShopID:           shopMongoID,
		Progress:         []model.RegistrationProgress{{Status: "initiated", EventTime: time.Now()}},
	}

	//now save this supplier log and if any error dont proceed return back with error
	if _, err := config.MongoDbDatabase.Collection(constants.SupplierRegistrationLogsModel).InsertOne(ctx, supplierLogsPayload); err != nil {
		return nil, err
	}

	//now make json of this payload
	jsonSupplierRegistrationPayload, err := json.Marshal(&supplierLogsPayload)

	if err != nil {
		return nil, err
	}

	//now make request to shop module to save supplier id in this shop mentioned by shopId
	resp, err := http.Post(constants.BaseURL+constants.ShopModuleAddSupplierURL, "application/json", bytes.NewBuffer(jsonSupplierRegistrationPayload))
	if err != nil {
		return nil, err
	}

	fmt.Println("response from shopmodule is ", resp)

	if resp.StatusCode != http.StatusCreated {
		//here it means it failed to add supplier to shop so simply return with error here
		return nil, errors.New("failed to add supplier to shop")
	}

	//now save this supplier details
	if _, err := sv.Repo.AddSupplier(ctx, supplierDetails); err != nil {
		return nil, err
	}

	//now as data is saved update supplierLogs with status completed otherwise  this supplier details will be deleted after some time
	updateQuery := bson.M{"$push": bson.M{"progress": model.RegistrationProgress{Status: "completed", EventTime: time.Now()}}}

	if _, err := config.MongoDbDatabase.Collection(constants.SupplierRegistrationLogsModel).
		UpdateOne(ctx, bson.M{"supplierId": supplierDetails.ID}, updateQuery); err != nil {
		return nil, err
	}

	//now return back nil,nil controller doesnt need mongo insertone result so send it also as nil
	return nil, nil
}

func (sv *Service) GetSupplierByName(ctx context.Context, supplierName string) (model.Supplier, error) {
	return sv.Repo.GetSupplierByName(ctx, supplierName)
}
