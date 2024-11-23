package handlers

import (
	"github.com/gofiber/fiber/v2"
	"merchant-bank-api/app/dto/response"
	"merchant-bank-api/app/repository"
)

type CustomerHandler struct {
	repo *repository.JSONRepository
}

func NewCustomerHandler(repo *repository.JSONRepository) *CustomerHandler {
	return &CustomerHandler{repo: repo}
}

func (h *CustomerHandler) GetAllCustomer(c *fiber.Ctx) error {
	var message = ""
	customers := h.repo.FindAllCustomers()

	if len(customers) > 0 {
		message = "Success retrieve all activity"
	} else {
		message = "Activities null"
	}

	var customersResponse []response.CustomerResponse

	for i := range customers {
		customerResponse := response.CustomerResponse{
			ID:       customers[i].ID,
			Username: customers[i].Username,
			Balance:  customers[i].Balance,
			Email:    customers[i].Email,
			Status:   customers[i].Status,
		}
		customersResponse = append(customersResponse, customerResponse)
	}

	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: message,
		Data:    customersResponse,
	})
}

func (h *CustomerHandler) GetCustomerById(c *fiber.Ctx) error {
	id := c.Params("id")

	customer, exists := h.repo.FindCustomerById(id)

	if exists == false {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "Customer not found",
		})
	}

	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Success retrieve customer",
		Data: response.CustomerResponse{
			ID:       customer.ID,
			Username: customer.Username,
			Balance:  customer.Balance,
			Email:    customer.Email,
			Status:   customer.Status,
		},
	})
}
