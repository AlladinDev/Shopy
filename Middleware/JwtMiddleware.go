// Package middleware provides middlewares for various uses
package middleware

import (
	"errors"
	"fmt"
	"net/http"

	utils "github.com/AlladinDev/Shopy/Utils"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
)

func JwtAuthMiddleware(c *fiber.Ctx) error {
	jwtCookie := c.Get("jwtToken")

	if jwtCookie == "" {
		return utils.ReturnAppError(errors.New("jwt header missing"), "Unauthorized", http.StatusBadRequest)
	}

	//now verify jwt token inside it
	tokenDetails, err := utils.VerifyJwt(jwtCookie)

	if err != nil {
		return utils.ReturnAppError(err, "Unauthorized", http.StatusBadRequest)
	}

	jwtClaims, ok := tokenDetails.(jwt.MapClaims)
	if !ok {
		return utils.ReturnAppError(errors.New("malformed token"), "Unauthorized", http.StatusBadRequest)
	}

	//set userId decoded from jwt into local storage of fibre for use in handlers
	userID, ok := jwtClaims["userId"].(string)
	if !ok {
		return utils.ReturnAppError(errors.New("malformed token"), "Unauthorized", http.StatusBadRequest)
	}
	c.Locals("userId", userID)

	//similary set  shopId in locals store also
	shopID, ok := jwtClaims["shopId"].(string)
	if !ok {
		return utils.ReturnAppError(errors.New("malformed token"), "Unauthorized", http.StatusBadRequest)
	}
	c.Locals("shopId", shopID)
	fmt.Println("shopid in jwt middleware is", c.Locals("shopId"), shopID)

	//now set userType also
	userType, ok := jwtClaims["userType"].(string)
	if !ok {
		return utils.ReturnAppError(errors.New("malformed token"), "Unauthorized", http.StatusBadRequest)
	}
	c.Locals("userType", userType)

	return c.Next()
}
