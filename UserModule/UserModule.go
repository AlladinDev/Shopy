// Package usermodule acts as central file for user module
package usermodule

import (
	config "UserService/Pkg/Config"
	worker "UserService/UserModule/Worker"
	"UserService/UserModule/controller"
	"UserService/UserModule/repository"
	service "UserService/UserModule/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

func InitialiseUserModule(appConfig *config.Config, router fiber.Router) {
	repository := repository.NewUserRepository(appConfig.MongoDatabase)
	service := service.CreateNewUserService(appConfig, repository)
	controller := controller.CreateNewUserController(appConfig, service)
	router.Post("/user/register", controller.RegisterUser)
	router.Post("/user/login", controller.LoginUser)
	router.Get("/user/Bulk", controller.GetBulkUsers)
	router.Get("/user/details", controller.GetUserByID)
	router.Post("/user/addshop", controller.AddShopToUser)

	//start workers for different tasks
	go worker.CheckFailedShopRegistrationEventsWorker(4*time.Minute, appConfig.RabbitMqConnection)
}
