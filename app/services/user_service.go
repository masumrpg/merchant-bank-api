package services

import (
	"merchant-bank-api/app/dto/response"
	"merchant-bank-api/app/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates and returns a new instance of UserService
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// GetAllUsers retrieves all users and returns a slice of UserResponse
func (s *UserService) GetAllUsers() ([]response.UserResponse, string) {
	users := s.repo.FindAllUser()
	var message string
	if len(users) > 0 {
		message = "Success retrieve all users"
	} else {
		message = "No users found"
	}

	// Convert users from models to response format
	var usersResponse []response.UserResponse
	for _, user := range users {
		usersResponse = append(usersResponse, response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Balance:  user.Balance,
			Email:    user.Email,
			Status:   user.Status,
		})
	}

	return usersResponse, message
}

// GetUserById retrieves a user by their ID and returns a UserResponse
func (s *UserService) GetUserById(id string) (response.UserResponse, bool) {
	user, exists := s.repo.FindUserById(id)
	if !exists {
		return response.UserResponse{}, false
	}

	// Return user details in the response format
	return response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Balance:  user.Balance,
		Email:    user.Email,
		Status:   user.Status,
	}, true
}
