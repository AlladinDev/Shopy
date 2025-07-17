// Package shopmodule acts as central file for shop module
package shopmodule

import (
	"time"

	middleware "github.com/AlladinDev/Shopy/Middleware"
	service "github.com/AlladinDev/Shopy/ShopModule/Service"
	worker "github.com/AlladinDev/Shopy/ShopModule/Worker"

	constants "github.com/AlladinDev/Shopy/Constants"
	config "github.com/AlladinDev/Shopy/Pkg/Config"
	controller "github.com/AlladinDev/Shopy/ShopModule/Controller"
	repository "github.com/AlladinDev/Shopy/ShopModule/Repository"

	"github.com/gofiber/fiber/v2"
)

func InitialiseShopModule(appConfig *config.Config, router fiber.Router) {
	repository := repository.ReturnNewShopRepository(appConfig.MongoDatabase)
	service := service.ReturnNewShopService(repository)
	controllers := controller.GetNewShopController(service)

	//now define routes
	router.Post("shop/register", middleware.JwtAuthMiddleware, middleware.RoleGuards([]string{constants.UserTypeUser}), controllers.RegisterShop)
	router.Post(constants.ShopModuleAddSupplierURL, controllers.RegisterSupplier)
	router.Get("shop/bulk", controllers.GetAllShops)

	//start the worker for checking failed shop registration events
	go worker.CheckFailedShopRegistrationEventsWorker(5*time.Minute, appConfig.RabbitMqConnection)
	go worker.CheckFailedSupplierRegistrationLogs()
}
