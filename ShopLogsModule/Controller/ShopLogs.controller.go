// Package controller provides various handlers for shop logs module
package controller

import (
	interfaces "UserService/ShopLogsModule/Interfaces"
	models "UserService/ShopLogsModule/Models"
	utils "UserService/Utils"
	"context"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	Service interfaces.IService
}

// ReturnNewController this function will return new controller with all handlers for interacting with shopLogs module
func ReturnNewController(service interfaces.IService) *Controller {
	return &Controller{
		Service: service,
	}
}

func (sc *Controller) AddLog(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 4*time.Second)
	defer cancel()

	//now extract details bind them
	var shopLogDetails models.ShopRegistrationLogs
	if err := c.BodyParser(&shopLogDetails); err != nil {
		return utils.ReturnAppError(err, "failed to parse json", http.StatusBadRequest)
	}

	//now do validation
	validator := validator.New()
	if err := validator.Struct(&shopLogDetails); err != nil {
		return utils.ReturnAppError(err, "Invalid Details", http.StatusBadRequest)
	}

	_, err := sc.Service.AddLog(ctx, shopLogDetails)
	if err != nil {
		return utils.ReturnAppError(err, "Failed to add shop", http.StatusInternalServerError)
	}

	return utils.AppSuccess(c, "Successfully Added ShopLog", nil, http.StatusCreated)
}

func (sc *Controller) UpdateLog(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 4*time.Second)
	defer cancel()
	//first parse details
	var shopLogDetails models.ShopRegistrationLogs

	if err := c.BodyParser(&shopLogDetails); err != nil {
		return err
	}

	//now validate details
	validator := validator.New()
	if err := validator.Struct(&shopLogDetails); err != nil {
		return err
	}

	_, err := sc.Service.UpdateLog(ctx, shopLogDetails)

	if err != nil {
		return err
	}

	return utils.AppSuccess(c, "Successfully Updated Shop Log Details", nil, http.StatusCreated)
}

func (sc *Controller) GetAllShopLogs(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 4*time.Second)
	defer cancel()
	logs, err := sc.Service.GetAllShopLogs(ctx)
	if err != nil {
		return err
	}

	return utils.AppSuccess(c, "Successfull Fetched Shop Logs", logs, http.StatusOK)
}
