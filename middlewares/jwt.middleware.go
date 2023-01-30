package middlewares

import (
	"errors"

	"github.com/Sebas3270/calendar-app-backend/helpers"
	"github.com/gofiber/fiber/v2"
)

func GetTokenInfo(c *fiber.Ctx) (string, string, error) {
	headers := c.GetReqHeaders()
	token, ok := headers["X-Token"]

	if ok == false {
		return "", "", errors.New("Not token in the request")
	}

	userId, userName, err := helpers.ValidateJwt(token)

	return userId, userName, err
}

func ValidateJwt(c *fiber.Ctx) error {

	_, _, err := GetTokenInfo(c)

	if err != nil {
		print(err)
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	return c.Next()
}
