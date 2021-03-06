package middlewares

import (
	"server/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type RequestHeader struct {
	Authorization string `reqHeader:"Authorization"`
}

func MiddlewareUser(c *fiber.Ctx) error {
	var headerData RequestHeader

	if err := c.ReqHeaderParser(&headerData); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	tokenString := headerData.Authorization
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(utils.GetENV("ACCESS_KEY_SECRET")), nil
	})

	if err != nil {
		if err.Error() == "token contains an invalid number of segments" {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "token not exist",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	c.Locals("user_id", claims["id"])
	c.Locals("user_role", claims["role"])

	return c.Next()
}
