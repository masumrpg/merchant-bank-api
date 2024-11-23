package handlers

import (
	"math/rand"
	"merchant-bank-api/app/dto/request"
	"merchant-bank-api/app/dto/response"
	"merchant-bank-api/app/models"
	"merchant-bank-api/app/repository"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	repo *repository.JSONRepository
}

func NewPaymentHandler(repo *repository.JSONRepository) *PaymentHandler {
	return &PaymentHandler{repo: repo}
}

func (h *PaymentHandler) ProcessPayment(c *fiber.Ctx) error {
	var paymentRequest = request.PaymentRequest

	if err := c.BodyParser(&paymentRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	fromUsername := c.Locals("username").(string)

	// Check if recipient exists
	toUser, exists := h.repo.FindCustomerByUsername(paymentRequest.ToUsername)
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "User not found",
		})
	}

	// Get sender details
	fromUser, exists := h.repo.FindCustomerByUsername(fromUsername)
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "User not found",
		})
	}

	// Validate balance
	if fromUser.Balance < paymentRequest.Amount {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "Balance not enough",
		})
	}

	// Process payment
	fromUser.Balance -= paymentRequest.Amount
	toUser.Balance += paymentRequest.Amount

	// Save updated balance
	if err := h.repo.UpdateCustomerBalance(fromUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Message: "Error updating user sender balance",
			Error:   err.Error(),
		})
	}
	if err := h.repo.UpdateCustomerBalance(toUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Message: "Error updating user receiver balance",
			Error:   err.Error(),
		})
	}

	// Log transaction
	transaction := models.Transaction{
		IDTransaction: strconv.Itoa(rand.Int()),
		Type:          "PAYMENT",
		FromUser:      fromUsername,
		ToUser:        paymentRequest.ToUsername,
		Amount:        paymentRequest.Amount,
		Details:       paymentRequest.Details,
		Timestamp:     time.Now(),
	}
	h.repo.SavePaymentHistory(transaction)

	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Payment processed successfully",
		Data:    transaction,
	})
}

func (h *PaymentHandler) GetAllPayment(c *fiber.Ctx) error {
	var message = ""
	payment := h.repo.FindAllPayment()
	if len(payment) > 0 {
		message = "Success retrieve all payment"
	} else {
		message = "Payment null"
	}

	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: message,
		Data:    payment,
	})
}

func (h *PaymentHandler) GetPaymentById(c *fiber.Ctx) error {
	id := c.Params("id")
	payment, exists := h.repo.FindPaymentById(id)
	if exists == false {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "Payment not found",
		})
	}

	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Success retrieve payment",
		Data:    payment,
	})
}
