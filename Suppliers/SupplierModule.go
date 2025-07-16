// Package suppliermodule provides function to initialse supplier module
package suppliermodule

import (
	"time"

	constants "github.com/AlladinDev/Shopy/Constants"
	middleware "github.com/AlladinDev/Shopy/Middleware"
	config "github.com/AlladinDev/Shopy/Pkg/Config"
	controller "github.com/AlladinDev/Shopy/Suppliers/Controller"
	repository "github.com/AlladinDev/Shopy/Suppliers/Repository"
	service "github.com/AlladinDev/Shopy/Suppliers/Service"
	worker "github.com/AlladinDev/Shopy/Suppliers/Worker"

	"github.com/gofiber/fiber/v2"
)

func InitialiseSupplierModule(router fiber.Router) {
	repository := repository.ReturnNewRepository(config.MongoDbDatabase)
	service := service.ReturnNewService(repository)
	controller := controller.ReturnNewController(service)

	router.Post("/supplier/register", middleware.JwtAuthMiddleware, middleware.RoleGuards([]string{constants.UserTypeUser}), controller.RegisterNewSupplier)
	router.Get("/supplier/bulk", middleware.JwtAuthMiddleware, middleware.RoleGuards([]string{constants.RoleUserAndAdmin}), controller.GetAllSuppliers)
	router.Get("/supplier/details", middleware.JwtAuthMiddleware, middleware.RoleGuards([]string{constants.RoleUserAndAdmin}), controller.GetSupplierByName)

	//start workers for different background tasks
	go worker.WCheckFailedSupplierRegistrations(6 * time.Second)
}
