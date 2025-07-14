// Package service provides various functions for handling business logic of shop module
package service

import (
	constants "UserService/Constants"
	contracts "UserService/Contracts"
	interfaces "UserService/ShopModule/Interfaces"
	models "UserService/ShopModule/Models"
	model "UserService/Suppliers/Model"
	utils "UserService/Utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShopService struct {
	Repository interfaces.IShopRepository
}

// this will make sure all functions of shop interface are implemented by this ShopService struct
var _ interfaces.IShopService = (*ShopService)(nil)

// ReturnNewShopService returns a new shop service with shop repository as dependency
func ReturnNewShopService(shopRepository interfaces.IShopRepository) *ShopService {
	return &ShopService{
		Repository: shopRepository,
	}
}

// GetAllShops func to get all shops
func (sv *ShopService) GetAllShops(ctx context.Context) ([]models.Shop, error) {
	return sv.Repository.GetAllShops(ctx)
}

// GetShopByName func to get shop by name
func (sv *ShopService) GetShopByName(ctx context.Context, shopName string) (models.Shop, error) {
	return sv.Repository.GetShopByName(ctx, shopName)
}

// RegisterShop func to register shop
func (sv *ShopService) RegisterShop(ctx context.Context, shopDetails models.Shop, userID string) (*mongo.InsertOneResult, error) {

	//convert userId string to mongodb id
	userMongoDBID, idErr := primitive.ObjectIDFromHex(userID)

	if idErr != nil {
		return nil, idErr
	}

	//first check if this user has has any shop registered because as of now only one shop per user is allowed
	_, err := sv.Repository.GetShopByUserID(ctx, userMongoDBID)

	//error is nill means shop is present otherwise if shop is not present it will return errNoDocuments error
	if err == nil {
		return nil, utils.ReturnAppError(errors.New("shop already exists for this user"), "Forbidden Shop Already Created", http.StatusBadRequest)
	}

	//assign userid to owner field in shopdetails
	shopDetails.Owner = userMongoDBID

	//now create mongodb id for shop
	shopDetails.ID = primitive.NewObjectID()

	//now store event that shop registration has been started in ShopLogs collection for that create payload
	shopLogsPayload := contracts.ShopRegistrationLogs{
		RegistrationDate: time.Now(),
		UserID:           shopDetails.Owner,
		ShopID:           shopDetails.ID,
		//it means under two minutes its progress should be updated to completed otherwise workers will pick the shoplog and emit its event for deletion
		//workers will think that its expiry time has gone and it is a failed event
		ExpiryTime: time.Now().Add(2 * time.Minute),
		IsShopRegistrationFailureNotificationSent: false,
		Progress: []contracts.ShopRegistrationProgress{
			{
				Status:      "initiated",
				EventTime:   time.Now(),
				HandlerFunc: "RegisterShop of Shops module",
			},
		},
	}

	//convert this payload into json
	jsonShopLogsData, err := json.Marshal(&shopLogsPayload)

	if err != nil {
		return nil, err
	}

	//now make api call to shop logs module
	resp, httpErr := http.Post(constants.BaseURL+"/shoplogs/add", "application/json", bytes.NewBuffer(jsonShopLogsData))

	if httpErr != nil {
		return nil, httpErr
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, utils.ReturnAppError(errors.New("unexpected error occurred"), "failed to create shop", http.StatusInternalServerError)
	}

	//now make call to usermodule so that this shopid also gets stored in user document
	resp, apiErr := http.Post(constants.BaseURL+"/user/addshop", "application/json", bytes.NewBuffer(jsonShopLogsData))

	if apiErr != nil {
		return nil, apiErr
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, utils.ReturnAppError(errors.New("failed to update user with shopid"), "Shop Registration Failed", http.StatusInternalServerError)
	}

	//now store this shop by calling repo method
	mongoRes, err := sv.Repository.CreateShop(ctx, shopDetails)

	if err != nil {
		return nil, err
	}

	//now again send http request to shoplogs module to update its status from initiated to completed
	shopLogsPayload.Progress = []contracts.ShopRegistrationProgress{
		{
			Status:      "completed",
			EventTime:   time.Now(),
			HandlerFunc: "Register shop service func of shop module",
		},
	}

	//convert shoplogspayload to json and send request
	jsonPayload, err := json.Marshal(&shopLogsPayload)
	if err != nil {
		return nil, err
	}

	//now make api call so that another event gets stored in shoplogs this time of status completed representing that shop registration was successfull
	req, err := http.NewRequest(http.MethodPatch, constants.BaseURL+"/shoplogs/update", bytes.NewBuffer(jsonPayload))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	//now send reponse back to controller
	return mongoRes, nil
}

func (sv *ShopService) AddSupplierToShop(ctx context.Context, supplierDetails model.SupplierRegistrationLogs) (*mongo.UpdateResult, error) {
	//first convert ids of shop and supplierid to mongodb id format because they will be in string format
	var err error
	supplierDetails.ShopID, err = primitive.ObjectIDFromHex(supplierDetails.ShopID.Hex())
	if err != nil {
		return nil, err
	}

	supplierDetails.SupplierID, err = primitive.ObjectIDFromHex(supplierDetails.SupplierID.Hex())
	if err != nil {
		return nil, err
	}

	//now as both ids have been converted into mongodb id format call repo function to add supplier to this shop
	return sv.Repository.AddSupplier(ctx, supplierDetails)
}
