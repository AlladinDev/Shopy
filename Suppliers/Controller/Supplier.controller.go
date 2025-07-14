// Package controller provides various handler functions for supplier module
package controller

import (
	config "UserService/Pkg/Config"
	interfaces "UserService/Suppliers/Interfaces"
	model "UserService/Suppliers/Model"
	utils "UserService/Utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	Sv interfaces.IService
}

func ReturnNewController(service interfaces.IService) *Controller {
	return &Controller{
		Sv: service,
	}
}

func (sc *Controller) RegisterNewSupplier(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	//for testing get shopId from query then in future get it from jwt token decoded
	shopID := c.Query("shopId")

	if shopID == "" {
		return utils.ReturnAppError(errors.New("shopid is required"), "ShopId is Required", http.StatusBadRequest)
	}

	//parse the details
	var supplierDetails model.Supplier
	if err := c.BodyParser(&supplierDetails); err != nil {
		return err
	}

	fmt.Printf("%v", supplierDetails)

	//now validate details
	if err := config.Validator.Struct(&supplierDetails); err != nil {
		fmt.Println(config.Validator)
		return err
	}

	//now send details to service layer function
	if _, err := sc.Sv.RegisterSupplier(ctx, shopID, supplierDetails); err != nil {
		return err
	}

	return utils.AppSuccess(c, "Supplier Registered Successfully", nil, http.StatusCreated)

}

func (sc *Controller) GetAllSuppliers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	//now get page ,limit from query
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")

	//now convert page,limit string into int as page ,limit comes in string format from url
	pageNumber, err := strconv.Atoi(page)

	if err != nil {
		return err
	}

	pageLimit, err := strconv.Atoi(limit)

	if err != nil {
		return err
	}

	suppliers, err := sc.Sv.GetAllSuppliers(ctx, pageNumber, pageLimit)

	if err != nil {
		return err
	}

	return utils.AppSuccess(c, "Successfully Fetched Suppliers", suppliers, http.StatusOK)
}

func (sc *Controller) GetSupplierByName(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 3*time.Second)
	defer cancel()

	//get name from query
	supplierName := c.Query("supplierName")

	if supplierName == "" {
		return utils.ReturnAppError(errors.New("supplier name is required"), "Bad Request", http.StatusBadRequest)
	}

	supplier, err := sc.Sv.GetSupplierByName(ctx, supplierName)

	if err != nil {
		return err
	}

	return utils.AppSuccess(c, "Successfully Fetched Supplier", supplier, http.StatusOK)
}
