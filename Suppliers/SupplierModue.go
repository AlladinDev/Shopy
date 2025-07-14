// Package suppliermodule provides function to initialse supplier module
package suppliermodule

import (
	config "UserService/Pkg/Config"
	controller "UserService/Suppliers/Controller"
	repository "UserService/Suppliers/Repository"
	service "UserService/Suppliers/Service"
	worker "UserService/Suppliers/Worker"
	"time"

	"github.com/gofiber/fiber/v2"
)

func InitialiseSupplierModule(router fiber.Router) {
	repository := repository.ReturnNewRepository(config.MongoDbDatabase)
	service := service.ReturnNewService(repository)
	controller := controller.ReturnNewController(service)

	router.Post("/supplier/register", controller.RegisterNewSupplier)
	router.Get("/supplier/bulk", controller.GetAllSuppliers)
	router.Get("/supplier/details", controller.GetSupplierByName)

	//start workers for different background tasks
	go worker.WCheckFailedSupplierRegistrations(6 * time.Second)
}
