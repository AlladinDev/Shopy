package utils

import "github.com/gofiber/fiber/v2"

type Success struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func AppSuccess(c *fiber.Ctx, msg string, data any, statusCode int) error {
	payload := Success{
		StatusCode: statusCode,
		Message:    msg,
		Data:       data,
	}
	return c.Status(statusCode).JSON(payload)
}
