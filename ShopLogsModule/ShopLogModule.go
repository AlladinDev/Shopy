// Package shoplogsmodule serves as entry file for shoplogs module
package shoplogsmodule

import (
	config "UserService/Pkg/Config"
	controller "UserService/ShopLogsModule/Controller"
	repository "UserService/ShopLogsModule/Repository"
	service "UserService/ShopLogsModule/Service"
	worker "UserService/ShopLogsModule/Worker"
	"time"

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
