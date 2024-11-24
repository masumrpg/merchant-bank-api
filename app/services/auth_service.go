package services

import (
	"errors"
	"math/rand"
	"merchant-bank-api/app/config"
	"merchant-bank-api/app/dto/request"
	"merchant-bank-api/app/dto/response"
	"merchant-bank-api/app/models"
	"merchant-bank-api/app/repository"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo     *repository.UserRepository
	activityRepo *repository.HistoryRepository
}

func NewAuthService(userRepo *repository.UserRepository, activityRepo *repository.HistoryRepository) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		activityRepo: activityRepo,
	}
}

func (s *AuthService) RegisterUser(req *request.RegisterRequest) (*response.RegisterResponse, error) {
	// Validate req
	if req.Username == "" || req.Password == "" || req.Email == "" {
		return nil, errors.New("username, email, and password are required")
	}

	// Check username and email
	if _, exists := s.userRepo.FindUserByUsername(req.Username); exists {
		return nil, errors.New("username is already taken")
	}
	if _, exists := s.userRepo.FindUserByEmail(req.Email); exists {
		return nil, errors.New("email is already taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Save to data
	user := &models.User{
		ID:       strconv.Itoa(rand.Int()),
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Balance:  req.Balance,
		Status:   "active",
	}
	if err := s.userRepo.SaveUser(user); err != nil {
		return nil, errors.New("failed to save user")
	}

	// Save to log
	activity := models.Activity{
		ID:        user.ID,
		Type:      "REGISTER",
		Username:  user.Username,
		Details:   "Successful registration",
		Timestamp: time.Now(),
	}
	s.activityRepo.SaveActivityHistory(activity)

	// Return response
	return &response.RegisterResponse{
		User: response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Balance:  user.Balance,
			Status:   user.Status,
		},
	}, nil
}

func (s *AuthService) LoginUser(req *request.LoginRequest) (*response.LoginResponse, error) {
	// Validate req
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	// Find user
	user, exists := s.userRepo.FindUserByUsername(req.Username)
	if !exists {
		return nil, errors.New("invalid credentials")
	}
	if user.Status != "active" {
		return nil, errors.New("account is inactive")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// If failed save to history
		s.activityRepo.SaveActivityHistory(models.Activity{
			Type:      "FAILED_LOGIN",
			Username:  req.Username,
			Details:   "Invalid password attempt",
			Timestamp: time.Now(),
		})
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.LoadConfig().JWTSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Save activity
	s.activityRepo.SaveActivityHistory(models.Activity{
		ID:        user.ID,
		Type:      "LOGIN",
		Username:  user.Username,
		Details:   "Successful login",
		Timestamp: time.Now(),
		Token:     tokenString,
		ExpiresIn: expirationTime,
	})

	// Return response
	return &response.LoginResponse{
		User: response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Balance:  user.Balance,
			Status:   user.Status,
		},
		Token:     tokenString,
		ExpiresIn: expirationTime.Unix() - time.Now().Unix(),
	}, nil
}

func (s *AuthService) LogoutUser(token string) error {
	// Delete log
	err := s.activityRepo.DeleteLoggedByToken(token)
	if err != nil {
		return err
	}
	return nil
}
