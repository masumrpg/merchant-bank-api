package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"merchant-bank-api/app/config"
	"merchant-bank-api/app/dto/response"
	"merchant-bank-api/app/repository"
	"time"
)

type JWTHandler struct {
	userRepo    *repository.UserRepository
	historyRepo *repository.HistoryRepository
}

func NewJWTHandler(userRepo *repository.UserRepository, historyRepo *repository.HistoryRepository) *JWTHandler {
	return &JWTHandler{
		userRepo:    userRepo,
		historyRepo: historyRepo,
	}
}

func JWTMiddleware(repo *JWTHandler) fiber.Handler {

	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Message: "Unauthorized",
				Error:   "Missing authorization token",
			})
		}

		// Validate blocked token in log
		_, exists := repo.historyRepo.FindLoggedByToken(tokenString)
		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
				Message: "You are logged out!",
				Error:   "The token is logged out and blocked!",
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
		_, existsByUsername := repo.userRepo.FindUserByUsername(username)
		if !existsByUsername {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Delete activities if expired
		go cleanUpExpiredActivities(repo.historyRepo)

		// Set username to context
		c.Locals("username", username)
		return c.Next()
	}
}

// cleanUpExpiredActivities for clean token if expired
func cleanUpExpiredActivities(r *repository.HistoryRepository) {
	activities := r.FindAllActivities()

	for _, activity := range activities {
		if activity.ExpiresIn.Before(time.Now()) {
			err := r.DeleteActivity(activity.ID)
			if err != nil {
				log.Println("Error deleting activity:", err)
				return
			}
		}
	}
}
