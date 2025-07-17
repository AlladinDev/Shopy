// Package usermodule acts as central file for user module
package usermodule

import (
	"time"

	constants "github.com/AlladinDev/Shopy/Constants"
	controllers "github.com/AlladinDev/Shopy/UserModule/controller"
	"github.com/AlladinDev/Shopy/UserModule/repository"
	service "github.com/AlladinDev/Shopy/UserModule/service"

	middleware "github.com/AlladinDev/Shopy/Middleware"
	config "github.com/AlladinDev/Shopy/Pkg/Config"
	worker "github.com/AlladinDev/Shopy/UserModule/Worker"

	"github.com/gofiber/fiber/v2"
)

func InitialiseUserModule(appConfig *config.Config, router fiber.Router) {
	repository := repository.NewUserRepository(appConfig.MongoDatabase)
	service := service.CreateNewUserService(appConfig, repository)
	controller := controllers.CreateNewUserController(appConfig, service)
	router.Post("/user/register", controller.RegisterUser)
	router.Post("/user/login", controller.LoginUser)
	router.Get("/user/Bulk", middleware.JwtAuthMiddleware, middleware.RoleGuards([]string{constants.RoleUserAndAdmin}), controller.GetBulkUsers)
	router.Get("/user/details", middleware.JwtAuthMiddleware, middleware.RoleGuards([]string{constants.RoleUserAndAdmin}), controller.GetUserByID)
	router.Post("/user/addshop", controller.AddShopToUser)

	//start workers for different tasks
	go worker.CheckFailedShopRegistrationEventsWorker(4*time.Minute, appConfig.RabbitMqConnection)
}
