// Package constants define various constants to be used throughout the app
package constants

const (
	UserTypeAdmin                      = "admin"
	UserTypeUser                       = "user"
	DatabaseName                       = "Restaurant"
	AdminCollection                    = "Admin"
	UserCollection                     = "User"
	RabbitMqShopRegistrationExchange   = "Shop_Registration_Exchange"
	ShopCollection                     = "Shop"
	ShopLogsCollection                 = "ShopLogs"
	SupplierModel                      = "Supplier"
	SupplierRegistrationLogsModel      = "SupplierRegistrationLogs"
	ShopRegistrationQueueForUserModule = "ShopRegistrationQueue_UserModule"
	ShopRegistrationQueueForShopModule = "ShopRegistrationQueue_ShopModule"
	RabbiqMQSupplierExchange           = "SupplierExchange"
	SupplierRegistrationFailedQueue    = "SupplierRegistrationFailed"
	NameOfDLXForSupplier               = "DLXSupplierExchange"
)

// these are the urls of the modules such as shop module user module for inter module communication
const (
	BaseURL                  = "http://localhost:3000/v1"
	ShopModuleAddSupplierURL = "/shop/supplier/add"
)
