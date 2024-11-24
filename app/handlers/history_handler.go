package handlers

import (
	"github.com/gofiber/fiber/v2"
	"merchant-bank-api/app/dto/response"
	"merchant-bank-api/app/models"
	"merchant-bank-api/app/services"
)

type HistoryHandler struct {
	service *services.HistoryService
}

// NewHistoryHandler creates and returns a new instance of HistoryHandler
func NewHistoryHandler(service *services.HistoryService) *HistoryHandler {
	return &HistoryHandler{service: service}
}

// GetAllActivities handles the retrieval of all activities
func (h *HistoryHandler) GetAllActivities(c *fiber.Ctx) error {
	responseData, message := h.service.GetAllActivities()
	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: message,
		Data:    responseData,
	})
}

// GetActivitiesByUsername handles the retrieval of activities by username
func (h *HistoryHandler) GetActivitiesByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	responseData, message := h.service.GetActivitiesByUsername(username)
	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: message,
		Data:    responseData,
	})
}

// SaveActivity handles saving a new activity
func (h *HistoryHandler) SaveActivity(c *fiber.Ctx) error {
	var activity models.Activity
	if err := c.BodyParser(&activity); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "Invalid input",
		})
	}

	h.service.SaveActivity(activity)

	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Activity saved successfully",
	})
}
