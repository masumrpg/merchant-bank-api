package handlers

import (
	"github.com/gofiber/fiber/v2"
	"merchant-bank-api/app/dto/request"
	"merchant-bank-api/app/dto/response"
	"merchant-bank-api/app/services"
)

type PaymentHandler struct {
	service *services.PaymentService
}

// NewPaymentHandler creates and returns a new instance of PaymentHandler
func NewPaymentHandler(service *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

// ProcessPayment handles the payment request from the client
func (h *PaymentHandler) ProcessPayment(c *fiber.Ctx) error {
	// Parse the payment request from the request body
	var paymentRequest request.PaymentRequest
	if err := c.BodyParser(&paymentRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	fromUsername := c.Locals("username").(string)

	// Call the service layer to process the payment
	transaction, err := h.service.ProcessPayment(paymentRequest, fromUsername)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Return success response
	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Payment processed successfully",
		Data:    transaction,
	})
}

// GetAllPayment handles the request to retrieve all payment transactions
func (h *PaymentHandler) GetAllPayment(c *fiber.Ctx) error {
	// Call the service layer to get all payments
	payments, err := h.service.GetAllPayments()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Message: "Error retrieving payments",
		})
	}

	// Return the response with the list of payments
	message := "Success retrieve all payments"
	if len(payments) == 0 {
		message = "No payments found"
	}

	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: message,
		Data:    payments,
	})
}

// GetPaymentById handles the request to retrieve a payment by its ID
func (h *PaymentHandler) GetPaymentById(c *fiber.Ctx) error {
	// Retrieve the payment ID from URL parameters
	id := c.Params("id")

	// Call the service layer to retrieve the payment
	payment, exists := h.service.GetPaymentById(id)
	if !exists {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "Payment not found",
		})
	}

	// Return the payment details
	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Payment retrieved successfully",
		Data:    payment,
	})
}
