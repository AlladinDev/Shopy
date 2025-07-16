package middleware

import (
	"errors"
	"net/http"

	utils "github.com/AlladinDev/Shopy/Utils"
	"github.com/gofiber/fiber/v2"
)

func RoleGuards(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//get userType from c.locals as it will be set by jwt middleware
		userType := c.Locals("userType").(string)
		if userType == "" {
			return utils.ReturnAppError(errors.New("userType is required"), "Authorisation Failed", http.StatusBadRequest)
		}

		//first check if roles array includes BothAdmin  it means both user and admin can access this endpoint and just jwt auth is required so pass this request
		if roles[0] == "Admin&User" {
			return c.Next()
		}

		//now check if userType is present in roles only then allow to go further
		for _, role := range roles {
			if userType == role {
				return c.Next()
			}
		}

		//if userType is not equal to any role mentioned in roles it means return error as this role is not supported
		return utils.ReturnAppError(errors.New("role is not allowed"), "UnAuthorised", http.StatusBadRequest)
	}
}
