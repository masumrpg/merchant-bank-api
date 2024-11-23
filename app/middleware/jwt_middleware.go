package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"merchant-bank-api/app/config"
	"merchant-bank-api/app/dto/response"
	"merchant-bank-api/app/repository"
)

func JWTMiddleware(repo *repository.JSONRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization token",
			})
		}

		// Validate token in log
		_, exists := repo.FindLoggedByToken(tokenString)
		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Message: "Auth is blocked please login again!",
			})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return config.LoadConfig().JWTSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		claims := token.Claims.(jwt.MapClaims)
		username := claims["username"].(string)

		// Validate user existence
		_, existsByUsername := repo.FindCustomerByUsername(username)
		if !existsByUsername {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		c.Locals("username", username)
		return c.Next()
	}
}
