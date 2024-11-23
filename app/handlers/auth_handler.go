package handlers

import (
	"math/rand"
	"merchant-bank-api/app/config"
	"merchant-bank-api/app/dto/request"
	"merchant-bank-api/app/dto/response"
	"merchant-bank-api/app/models"
	"merchant-bank-api/app/repository"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	repo *repository.JSONRepository
}

func NewAuthHandler(repo *repository.JSONRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var registerRequest request.RegisterRequest

	if err := c.BodyParser(&registerRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate required fields
	if registerRequest.Username == "" || registerRequest.Password == "" || registerRequest.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "Username, email and password are required",
		})
	}

	// Check if username already exists
	if _, exists := h.repo.FindCustomerByUsername(registerRequest.Username); exists {
		return c.Status(fiber.StatusConflict).JSON(response.ErrorResponse{
			Message: "Username is already taken",
		})
	}

	// Check if email already exists
	if _, exists := h.repo.FindCustomerByEmail(registerRequest.Email); exists {
		return c.Status(fiber.StatusConflict).JSON(response.ErrorResponse{
			Message: "Email is already taken",
		})
	}

	// Hash the password
	hashedPassword, err := hashPassword(registerRequest.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Message: "Error generating password",
			Error:   err.Error(),
		})
	}

	// Create customer
	customer := models.Customer{
		ID:       strconv.Itoa(rand.Int()),
		Username: registerRequest.Username,
		Password: hashedPassword,
		Email:    registerRequest.Email,
		Balance:  registerRequest.Balance,
		Status:   "active",
	}

	// Save the customer in repo
	if err := h.repo.SaveCustomer(&customer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Message: "Error saving customer",
			Error:   err.Error(),
		})
	}

	h.repo.SaveActivityHistory(models.Activity{
		ID:        customer.ID,
		Type:      "REGISTER",
		Username:  customer.Username,
		Details:   "Successful login",
		Timestamp: time.Now(),
	})

	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusCreated,
		Message: "Registration successful",
		Data: response.RegisterResponse{
			Customer: response.CustomerResponse{
				ID:       customer.ID,
				Username: customer.Username,
				Balance:  customer.Balance,
				Email:    customer.Email,
				Status:   customer.Status,
			},
		},
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var loginRequest request.LoginRequest

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	// Validate required fields
	if loginRequest.Username == "" || loginRequest.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse{
			Message: "Username and password are required",
		})
	}

	// Find customer
	customer, exists := h.repo.FindCustomerByUsername(loginRequest.Username)
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
			Message: "Invalid credentials",
		})
	}

	// Check if customer is active
	if customer.Status != "active" {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
			Message: "Account is inactive",
		})
	}

	// Verify password
	err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(loginRequest.Password))
	if err != nil {
		// Log failed login attempt
		h.repo.SaveActivityHistory(models.Activity{
			Type:      "FAILED_LOGIN",
			Username:  loginRequest.Username,
			Details:   "Invalid password attempt",
			Timestamp: time.Now(),
		})

		return c.Status(fiber.StatusUnauthorized).JSON(response.ErrorResponse{
			Message: "Invalid credentials",
		})
	}

	// Generate JWT Token
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := jwt.MapClaims{
		"user_id":  customer.ID,
		"username": customer.Username,
		"exp":      expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.LoadConfig().JWTSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse{
			Message: "Error generating token",
			Error:   err.Error(),
		})
	}

	// Log successful login
	h.repo.SaveActivityHistory(models.Activity{
		ID:        customer.ID,
		Type:      "LOGIN",
		Username:  customer.Username,
		Details:   "Successful login",
		Timestamp: time.Now(),
		Token:     tokenString,
		ExpiresIn: expirationTime,
	})

	return c.JSON(response.SuccessResponse{
		Status:  fiber.StatusOK,
		Message: "Successful login",
		Data: response.LoginResponse{
			Customer: response.CustomerResponse{
				ID:       customer.ID,
				Username: customer.Username,
				Balance:  customer.Balance,
				Email:    customer.Email,
				Status:   customer.Status,
			},
			Token:     tokenString,
			ExpiresIn: expirationTime.Unix() - time.Now().Unix(),
		},
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	token := c.Get("Authorization")

	// Delete log
	err := h.repo.DeleteLoggedByToken(token)
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

// hashPassword Helper function to hash password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
