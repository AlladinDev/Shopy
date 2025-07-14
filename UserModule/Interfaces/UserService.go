package interfaces

import (
	"context"

	contracts "github.com/AlladinDev/Shopy/Contracts"
	models "github.com/AlladinDev/Shopy/UserModule/Models"
	schemas "github.com/AlladinDev/Shopy/UserModule/Schemas"

	"go.mongodb.org/mongo-driver/mongo"
)

type IUserService interface {
	RegisterUser(ctx context.Context, userDetails models.User) (*mongo.InsertOneResult, error)
	LoginUser(ctx context.Context, loginDetails schemas.UserLoginDTO) (jwtToken string, err error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByPhoneNumber(ctx context.Context, mobileNumber int) (models.User, error)
	GetUserByID(ctx context.Context, userID string) (models.User, error)
	AddShop(ctx context.Context, shopDetails contracts.ShopRegistrationLogs) *mongo.SingleResult
}
