// Package shopmodule acts as central file for shop module
package shopmodule

import (
	constants "UserService/Constants"
	config "UserService/Pkg/Config"
	controller "UserService/ShopModule/Controller"
	repository "UserService/ShopModule/Repository"
	service "UserService/ShopModule/Service"
	worker "UserService/ShopModule/Worker"
	"time"

	"github.com/gofiber/fiber/v2"
)

func InitialiseShopModule(appConfig *config.Config, router fiber.Router) {
	repository := repository.ReturnNewShopRepository(appConfig.MongoDatabase)
	service := service.ReturnNewShopService(repository)
	controllers := controller.GetNewShopController(service)

	//now define routes
	router.Post("shop/register", controllers.RegisterShop)
	router.Post(constants.ShopModuleAddSupplierURL, controllers.RegisterSupplier)

	//start the worker for checking failed shop registration events
	go worker.CheckFailedShopRegistrationEventsWorker(5*time.Minute, appConfig.RabbitMqConnection)
	go worker.CheckFailedSupplierRegistrationLogs()
}
