package handlers

import (
	"github.com/gofiber/fiber/v2"
	"merchant-bank-api/app/dto/response"
	"merchant-bank-api/app/services"
)

type UserHandler struct {
	service *services.UserService
}

// NewUserHandler creates and returns a new instance of UserHandler
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetAllUser handles the request to retrieve all users
func (h *UserHandler) GetAllUser(c *fiber.Ctx) error {
	// Call the service layer to get all users
	usersResponse, message := h.service.GetAllUsers()

	// Return the response with the list of users
	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: message,
		Data:    usersResponse,
	})
}

// GetUserById handles the request to retrieve a user by their ID
func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	// Retrieve the user ID from URL parameters
	id := c.Params("id")

	// Call the service layer to get the user
	userResponse, exists := h.service.GetUserById(id)
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "User not found",
		})
	}

	// Return the user details
	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "User retrieved successfully",
		Data:    userResponse,
	})
}
