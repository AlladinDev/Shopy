// Package shoplogsmodule serves as entry file for shoplogs module
package shoplogsmodule

import (
	"time"

	config "github.com/AlladinDev/Shopy/Pkg/Config"
	controller "github.com/AlladinDev/Shopy/ShopLogsModule/Controller"
	repository "github.com/AlladinDev/Shopy/ShopLogsModule/Repository"
	service "github.com/AlladinDev/Shopy/ShopLogsModule/Service"
	worker "github.com/AlladinDev/Shopy/ShopLogsModule/Worker"

	"github.com/gofiber/fiber/v2"
)

func InitialiseShopLogsModule(appConfig *config.Config, router fiber.Router) {
	repository := repository.ReturnNewRepository(appConfig.MongoDatabase)
	service := service.ReturnNewService(repository)
	controller := controller.ReturnNewController(service)

	router.Post("/shoplogs/add", controller.AddLog)
	router.Patch("/shoplogs/update", controller.UpdateLog)
	router.Post("/shoplogs/bulk", controller.GetAllShopLogs)

	//spin worker for checking failed shop registrations
	go worker.ShopRegistrationFailureCheckingWorker(3*time.Minute, appConfig.RabbitMqConnection)
}
