package middleware

import (
	utils "UserService/Utils"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func JwtAuthMiddleware(c *fiber.Ctx) {
	jwtCookie := c.Get("jwtToken")

	if jwtCookie == "" {
		c.Status(http.StatusUnauthorized).JSON(utils.AppError{Msg: "Jwt Token missing", StatusCode: http.StatusUnauthorized})
		return
	}

	//now verify jwt token inside it
	tokenDetails, err := utils.VerifyJwt(jwtCookie)

	if err != nil {
		c.Status(http.StatusUnauthorized).JSON(utils.AppError{Msg: "Jwt Token missing", StatusCode: http.StatusUnauthorized})
		return
	}

	fmt.Println(tokenDetails)
	c.Next()

}
