// Package controller provides handlers for user module
package controller

import (
	"context"
	"errors"
	"net/http"
	"time"

	contracts "github.com/AlladinDev/Shopy/Contracts"
	config "github.com/AlladinDev/Shopy/Pkg/Config"
	interfaces "github.com/AlladinDev/Shopy/UserModule/Interfaces"
	models "github.com/AlladinDev/Shopy/UserModule/Models"
	schemas "github.com/AlladinDev/Shopy/UserModule/Schemas"
	utils "github.com/AlladinDev/Shopy/Utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	UserService interfaces.IUserService
	appConfig   *config.Config
}

func CreateNewUserController(appConfig *config.Config, sv interfaces.IUserService) *UserController {
	return &UserController{
		UserService: sv,
		appConfig:   appConfig,
	}
}

func (uc *UserController) RegisterUser(c *fiber.Ctx) error {
	var userDetails models.User
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()
	if err := c.BodyParser(&userDetails); err != nil {
		return err
	}

	_, err := uc.UserService.RegisterUser(ctx, userDetails)
	if err != nil {
		return err
	}

	return utils.AppSuccess(c, "User Registered Successfully", nil, http.StatusCreated)
}

func (uc *UserController) LoginUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 4*time.Second)
	defer cancel()

	//first bind data
	var loginDetails schemas.UserLoginDTO
	if err := c.BodyParser(&loginDetails); err != nil {
		return err
	}

	token, err := uc.UserService.LoginUser(ctx, loginDetails)
	if err != nil {
		return err
	}

	//prepare cookie
	cookie := &fiber.Cookie{
		Name:     "authToken",
		Value:    token,
		Path:     "/",
		Secure:   true,
		HTTPOnly: true,
		MaxAge:   3 * 24 * 60 * 60,
		Domain:   "",
		Expires:  time.Now().Add(3 * 24 * time.Minute),
	}

	c.Cookie(cookie)

	//now send login successfull message
	return utils.AppSuccess(c, "Login Successfull", nil, http.StatusOK)
}

func (uc *UserController) GetBulkUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 4*time.Second)
	defer cancel()

	//get all users
	users, err := uc.UserService.GetAllUsers(ctx)

	if err != nil {
		return err
	}

	return utils.AppSuccess(c, "Success", users, http.StatusOK)
}

func (uc *UserController) GetUserByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 4*time.Second)
	defer cancel()

	//get the id from locals as it will be set by jwt guard
	userID, ok := c.Locals("userId").(string)

	if !ok {
		return utils.ReturnAppError(errors.New("user id is missing"), "UserId missing", http.StatusBadRequest)
	}

	user, err := uc.UserService.GetUserByID(ctx, userID)

	if err != nil {
		return err
	}

	return utils.AppSuccess(c, "Successfully Fetched user", user, http.StatusOK)
}

func (uc *UserController) AddShopToUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 4*time.Second)
	defer cancel()

	//first parse response body
	var shopDetails contracts.ShopRegistrationLogs
	if err := c.BodyParser(&shopDetails); err != nil {
		return err
	}

	//now validate details
	validator := validator.New()
	if err := validator.Struct(&shopDetails); err != nil {
		return err
	}

	var err error

	//convert shopId which is in string format to mongodb format because data comes in json format and ids are in string form
	///so convert them into mongodb id format only then mongodb operations will be successfull because stringId is not equal to mongodbId  even if they are same
	shopDetails.ShopID, err = primitive.ObjectIDFromHex(shopDetails.ShopID.Hex())
	if err != nil {
		return err
	}

	shopDetails.UserID, err = primitive.ObjectIDFromHex(shopDetails.UserID.Hex())
	if err != nil {
		return err
	}

	//now call the service function
	_ = uc.UserService.AddShop(ctx, shopDetails)

	return utils.AppSuccess(c, "Shop Added To User Successfully", nil, http.StatusCreated)
}
