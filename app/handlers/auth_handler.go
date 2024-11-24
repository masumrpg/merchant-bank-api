package handlers

import (
	"merchant-bank-api/app/dto/request"
	"merchant-bank-api/app/dto/response"
	"merchant-bank-api/app/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var registerRequest request.RegisterRequest
	if err := c.BodyParser(&registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Registration service
	result, err := h.authService.RegisterUser(&registerRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "Successfully create user",
		Data:    result,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var loginRequest request.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Login service
	res, err := h.authService.LoginUser(&loginRequest)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Successfully login user",
		Data:    res,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	token := c.Get("Authorization")

	// Logout service
	err := h.authService.LogoutUser(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Message: "Error deleting token",
			Error:   err.Error(),
		})
	}

	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Successfully logged out",
		Data:    map[string]string{"username": username},
	})
}
