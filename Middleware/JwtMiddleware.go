// Package middleware provides middlewares for various uses
package middleware

import (
	"errors"
	"fmt"
	"net/http"

	utils "github.com/AlladinDev/Shopy/Utils"

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
		return utils.ReturnAppError(errors.New("malformed jwt"), "Unauthorized", http.StatusBadRequest)
	}

	fmt.Println("token detalils", tokenDetails)
	return c.Next()
}
