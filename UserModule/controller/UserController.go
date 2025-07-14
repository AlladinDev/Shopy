// Package controller provides handlers for user module
package controller

import (
	contracts "UserService/Contracts"
	config "UserService/Pkg/Config"
	interfaces "UserService/UserModule/Interfaces"
	models "UserService/UserModule/Models"
	schemas "UserService/UserModule/Schemas"
	utils "UserService/Utils"
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserController struct {
	userService interfaces.IUserService
	appConfig   *config.Config
}

func CreateNewUserController(appConfig *config.Config, userService interfaces.IUserService) *UserController {
	return &UserController{
		userService: userService,
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

	_, err := uc.userService.RegisterUser(ctx, userDetails)
	return err
}

func (uc *UserController) LoginUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 4*time.Second)
	defer cancel()

	//first bind data
	var loginDetails schemas.UserLoginDTO
	if err := c.BodyParser(&loginDetails); err != nil {
		return err
	}

	token, err := uc.userService.LoginUser(ctx, loginDetails)
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
	users, err := uc.userService.GetAllUsers(ctx)

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

	user, err := uc.userService.GetUserByID(ctx, userID)

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

	//convert ids which are in string format to mongodb format because data comes in json format and ids are in string form
	///so convert them into mongodb id format only then mongodb operations will be successfull because stringId is not equal to mongodbId  even if they are same
	shopID, err := primitive.ObjectIDFromHex(shopDetails.ShopID.Hex())
	if err != nil {
		return err
	}

	//overwrite shopId with its mongodb id format
	shopDetails.ShopID = shopID

	userID, err := primitive.ObjectIDFromHex(shopDetails.UserID.Hex())
	if err != nil {
		return err
	}

	//overwrite userid with its mongodb id format
	shopDetails.UserID = userID

	//now call the service function
	_ = uc.userService.AddShop(ctx, shopDetails)

	return utils.AppSuccess(c, "Shop Added To User Successfully", nil, http.StatusCreated)
}
