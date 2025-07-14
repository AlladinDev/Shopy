// Package controller provides various handlers for shop module
package controller

import (
	config "UserService/Pkg/Config"
	interfaces "UserService/ShopModule/Interfaces"
	models "UserService/ShopModule/Models"
	model "UserService/Suppliers/Model"
	utils "UserService/Utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ShopController struct {
	ShopService interfaces.IShopService
	AppConfig   *config.Config
}

// GetNewShopController func will get new shop controller with ShopService as dependency
func GetNewShopController(shopService interfaces.IShopService) *ShopController {
	return &ShopController{
		ShopService: shopService,
	}
}

// RegisterShop func to register shop
func (sc *ShopController) RegisterShop(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 4*time.Second)
	defer cancel()

	//first parse data
	var shopDetails models.Shop
	if err := c.BodyParser(&shopDetails); err != nil {
		return err
	}

	fmt.Print("hi passed there 1")

	//now validate shop details
	validator := validator.New()
	if err := validator.Struct(&shopDetails); err != nil {
		return err
	}

	fmt.Print("hi passed there")
	// //now extract userId from app context set by jwt middleware
	// userID, ok := c.Locals("userId").(string)

	// //if userid is not present return error
	// if !ok {
	// 	return utils.ReturnAppError(errors.New("userId is required for registering shop"), "UserId is missing", http.StatusBadRequest)
	// }

	//now call service function to register shop
	_, err := sc.ShopService.RegisterShop(ctx, shopDetails, shopDetails.Owner.Hex())

	if err != nil {
		return err
	}

	return utils.AppSuccess(c, "Registered Successfully", nil, http.StatusCreated)

}

// GetAllShops func to get all shops
func (sc *ShopController) GetAllShops(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 4*time.Second)
	defer cancel()
	shops, err := sc.ShopService.GetAllShops(ctx)

	if err != nil {
		return err
	}

	return utils.AppSuccess(c, "Fetched All Shops Successfully", shops, http.StatusOK)
}

// GetShopByName func to get shop by name
func (sc *ShopController) GetShopByName(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 4*time.Second)
	defer cancel()

	//get shopname by params
	shopName := c.Params("shopName")

	if shopName == "" {
		return utils.ReturnAppError(errors.New("shopname is required"), "ShopName is Required", http.StatusBadRequest)
	}

	shop, err := sc.ShopService.GetShopByName(ctx, shopName)
	if err != nil {
		return err
	}

	return utils.AppSuccess(c, "Fetched Shop", shop, http.StatusOK)
}

func (sc *ShopController) RegisterSupplier(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	//parse the body
	var supplierDetails model.SupplierRegistrationLogs
	if err := c.BodyParser(&supplierDetails); err != nil {
		return utils.ReturnAppError(err, "Body Parsing Error", http.StatusBadRequest)
	}

	//now do validation
	validator := validator.New()
	if err := validator.Struct(&supplierDetails); err != nil {
		return utils.ReturnAppError(err, "Validation Failed", http.StatusBadRequest)
	}

	//now call service function for further process
	if _, err := sc.ShopService.AddSupplierToShop(ctx, supplierDetails); err != nil {
		return utils.ReturnAppError(err, "Failed to add Supplier To Shop", http.StatusInternalServerError)
	}

	return utils.AppSuccess(c, "Successfully Added Supplier To Shop", nil, http.StatusCreated)
}
