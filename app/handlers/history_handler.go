package handlers

import (
	"github.com/gofiber/fiber/v2"
	"merchant-bank-api/app/dto/response"
)

func (h *AuthHandler) GetAllActivities(c *fiber.Ctx) error {
	var message = ""
	activities := h.repo.FindAllActivities()
	if len(activities) > 0 {
		message = "Success retrieve all activity"
	} else {
		message = "Activities null"
	}

	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: message,
		Data:    activities,
	})
}

func (h *AuthHandler) GetActivityByUsername(c *fiber.Ctx) error {
	var message = ""
	username := c.Params("username")
	activities := h.repo.FindActivitiesByUsername(username)

	if len(activities) > 0 {
		message = "Success retrieve all activity"
	} else {
		message = "Activities null"
	}

	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: message,
		Data:    activities,
	})
}
