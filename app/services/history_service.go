package services

import (
	"merchant-bank-api/app/dto/response"
	"merchant-bank-api/app/models"
	"merchant-bank-api/app/repository"
)

type HistoryService struct {
	repo *repository.HistoryRepository
}

// NewHistoryService creates a new instance of HistoryService
func NewHistoryService(repo *repository.HistoryRepository) *HistoryService {
	return &HistoryService{repo: repo}
}

// GetAllActivities retrieves all activities and returns a response message
func (s *HistoryService) GetAllActivities() (response.SuccessResponse, string) {
	activities := s.repo.FindAllActivities()
	var message string
	var responseData interface{}

	if len(activities) > 0 {
		message = "Success retrieve all activity"
		responseData = activities
	} else {
		message = "Activities null"
		responseData = []models.Activity{}
	}

	return response.SuccessResponse{
		Status:  200,
		Message: message,
		Data:    responseData,
	}, message
}

// GetActivitiesByUsername retrieves activities for a specific user and returns a response message
func (s *HistoryService) GetActivitiesByUsername(username string) (response.SuccessResponse, string) {
	activities := s.repo.FindActivitiesByUsername(username)
	var message string
	var responseData interface{}

	if len(activities) > 0 {
		message = "Success retrieve all activity"
		responseData = activities
	} else {
		message = "Activities null"
		responseData = []models.Activity{}
	}

	return response.SuccessResponse{
		Status:  200,
		Message: message,
		Data:    responseData,
	}, message
}

// SaveActivity saves a new activity
func (s *HistoryService) SaveActivity(activity models.Activity) {
	s.repo.SaveActivityHistory(activity)
}
