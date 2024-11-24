package services

import (
	"errors"
	"math/rand"
	"merchant-bank-api/app/dto/request"
	"merchant-bank-api/app/models"
	"merchant-bank-api/app/repository"
	"strconv"
	"time"
)

type PaymentService struct {
	transactionRepo *repository.TransactionRepository
	userRepo        *repository.UserRepository
}

// NewPaymentService creates and returns a new instance of PaymentService
func NewPaymentService(transactionRepo *repository.TransactionRepository, userRepo *repository.UserRepository) *PaymentService {
	return &PaymentService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

// ProcessPayment processes the payment transaction between two users
func (s *PaymentService) ProcessPayment(paymentRequest request.PaymentRequest, fromUsername string) (models.Transaction, error) {
	// Check if recipient exists
	toUser, exists := s.userRepo.FindUserByUsername(paymentRequest.ToUsername)
	if !exists {
		return models.Transaction{}, errors.New("user not found")
	}

	// Get sender details
	fromUser, exists := s.userRepo.FindUserByUsername(fromUsername)
	if !exists {
		return models.Transaction{}, errors.New("user not found")
	}

	// Validate balance
	if fromUser.Balance < paymentRequest.Amount {
		return models.Transaction{}, errors.New("insufficient balance")
	}

	// Process the payment by updating the balances
	fromUser.Balance -= paymentRequest.Amount
	toUser.Balance += paymentRequest.Amount

	// Update both sender and receiver balance
	if err := s.userRepo.UpdateUserBalance(fromUser); err != nil {
		return models.Transaction{}, err
	}
	if err := s.userRepo.UpdateUserBalance(toUser); err != nil {
		return models.Transaction{}, err
	}

	// Log the transaction
	transaction := models.Transaction{
		IDTransaction: strconv.Itoa(rand.Int()), // Generating a random ID for the transaction
		Type:          "PAYMENT",
		FromUser:      fromUsername,
		ToUser:        paymentRequest.ToUsername,
		Amount:        paymentRequest.Amount,
		Details:       paymentRequest.Details,
		Timestamp:     time.Now(),
	}
	s.transactionRepo.SavePaymentHistory(transaction)

	return transaction, nil
}

// GetAllPayments retrieves all payments from the database
func (s *PaymentService) GetAllPayments() ([]models.Transaction, error) {
	return s.transactionRepo.FindAllPayment(), nil
}

// GetPaymentById retrieves a payment by its ID
func (s *PaymentService) GetPaymentById(id string) (models.Transaction, bool) {
	return s.transactionRepo.FindPaymentById(id)
}
