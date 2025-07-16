// Package service provides service functions for user service module
package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	constants "github.com/AlladinDev/Shopy/Constants"
	contracts "github.com/AlladinDev/Shopy/Contracts"
	config "github.com/AlladinDev/Shopy/Pkg/Config"
	interfaces "github.com/AlladinDev/Shopy/UserModule/Interfaces"
	models "github.com/AlladinDev/Shopy/UserModule/Models"
	schemas "github.com/AlladinDev/Shopy/UserModule/Schemas"
	utils "github.com/AlladinDev/Shopy/Utils"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserServce struct {
	UserRepo  interfaces.IUserRepository
	AppConfig *config.Config
}

func CreateNewUserService(config *config.Config, userRepository interfaces.IUserRepository) *UserServce {
	return &UserServce{
		UserRepo:  userRepository,
		AppConfig: config,
	}
}

// now enforec that this github.com/AlladinDev/Shopy should implement all methods of Igithub.com/AlladinDev/Shopy interface
var _ interfaces.IUserService = (*UserServce)(nil)

func (sv *UserServce) RegisterUser(ctx context.Context, userDetails models.User) (*mongo.InsertOneResult, error) {
	//first check if this userPhoneNumber or email already exists and if yes return error
	//as phoneNumber and email should be unique
	filter := bson.D{
		{
			//this is match using or so if userEmail or phoneNumber matches it will return the record
			Key: "$or", Value: bson.A{
				bson.D{{Key: "phoneNumber", Value: userDetails.PhoneNumber}},
				bson.D{{Key: "email", Value: userDetails.Email}},
			},
		},
	}

	usersAlreadyPresent, err := sv.AppConfig.MongoDatabase.Collection(constants.UserCollection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, utils.ReturnAppError(err, "Registration Failed", http.StatusInternalServerError)
	}

	//it means user is already present with this phonenumber and email so return error
	if usersAlreadyPresent > 0 {
		return nil, utils.ReturnAppError(errors.New("user already exists with this phoneNumber or email"), "Registration Failed", http.StatusConflict)
	}

	//now set createdAt date
	userDetails.CreatedAt = time.Now()

	//now set default shop id
	userDetails.Shop = primitive.NilObjectID

	//now hash the password
	hash, hashingErr := bcrypt.GenerateFromPassword([]byte(userDetails.Password), 10)

	if hashingErr != nil {
		return nil, utils.ReturnAppError(hashingErr, "Registration Failed", http.StatusInternalServerError)
	}

	//now override plain password with hashed password
	userDetails.Password = string(hash)

	//now call the method of userRepo to register user
	return sv.UserRepo.RegisterUser(ctx, userDetails)
}

// GetAllUsers function to get all users
func (sv *UserServce) GetAllUsers(ctx context.Context) ([]models.User, error) {

	return sv.UserRepo.GetBulkUsers(ctx)

}

// GetUserByID func to get user by userid
func (sv *UserServce) GetUserByID(ctx context.Context, userID string) (models.User, error) {
	//first convert userid into mongo objectid
	userMongodbID, idErr := primitive.ObjectIDFromHex(userID)
	fmt.Print(userMongodbID, idErr)
	if idErr != nil {
		return models.User{}, idErr
	}

	//now call method user repo to get user
	return sv.UserRepo.GetUserByID(ctx, userMongodbID)
}

// LoginUser func to login user
func (sv *UserServce) LoginUser(ctx context.Context, loginDetails schemas.UserLoginDTO) (jwtToken string, err error) {
	user, err := sv.UserRepo.GetUserByEmail(ctx, loginDetails.Email)

	if err != nil {
		return "", nil
	}
	fmt.Println(user)

	payload := jwt.MapClaims{
		"userId":   user.ID,
		"userType": user.UserType,
		"shopId":   user.Shop,
		"iat":      time.Now(),
		"exp":      time.Now().Add(3 * 24 * time.Hour).Unix(),
	}

	//now generate jwt token
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", errors.New("unable to login")
	}

	token, err := utils.GenerateJwtToken(jwt.SigningMethodHS256, payload, secretKey)

	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByPhoneNumber func to get user by phone number
func (sv *UserServce) GetUserByPhoneNumber(ctx context.Context, mobileNumber int) (models.User, error) {
	//check if user exists or not using phoneNumber
	user, err := sv.UserRepo.GetUserByPhoneNumber(ctx, mobileNumber)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (sv *UserServce) AddShop(ctx context.Context, shopDetails contracts.ShopRegistrationLogs) *mongo.SingleResult {
	return sv.UserRepo.AddShop(ctx, shopDetails)
}
