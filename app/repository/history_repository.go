package repository

import (
	"encoding/json"
	"log"
	"merchant-bank-api/app/models"
	"os"
	"sync"
)

type HistoryRepository struct {
	mu sync.RWMutex
}

func NewHistoryRepository() *HistoryRepository {
	return &HistoryRepository{}
}

func (r *HistoryRepository) SaveActivityHistory(activity models.Activity) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Read existing activities
	data, err := os.ReadFile("data/activities.json")
	if err != nil {
		log.Println("Error reading activity file:", err)
		return
	}

	var activities []models.Activity
	errUnmarshal := json.Unmarshal(data, &activities)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
	}

	// Append new activity
	activities = append(activities, activity)

	// Write back to file
	updatedData, _ := json.MarshalIndent(activities, "", "  ")
	errWrite := os.WriteFile("data/activities.json", updatedData, 0644)
	if errWrite != nil {
		log.Fatalf("Error during Unmarshal: %v", errWrite)
	}
}

func (r *HistoryRepository) FindAllActivities() []models.Activity {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile("data/activities.json")
	if err != nil {
		log.Println("Error reading activities file:", err)
		return []models.Activity{}
	}

	var activities []models.Activity
	errUnmarshal := json.Unmarshal(data, &activities)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
		return []models.Activity{}
	}
	return activities
}

func (r *HistoryRepository) FindActivitiesByUsername(username string) []models.Activity {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile("data/activities.json")
	if err != nil {
		log.Println("Error reading activities file:", err)
		return []models.Activity{}
	}

	var activities []models.Activity
	errUnmarshal := json.Unmarshal(data, &activities)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
		return []models.Activity{}
	}

	var userActivities []models.Activity
	for _, activity := range activities {
		if activity.Username == username {
			userActivities = append(userActivities, activity)
		}
	}

	return userActivities
}

func (r *HistoryRepository) DeleteLoggedByToken(tokenString string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Read file
	data, err := os.ReadFile("data/activities.json")
	if err != nil {
		log.Println("Error reading activities file:", err)
		return err
	}

	var activities []models.Activity
	errUnmarshal := json.Unmarshal(data, &activities)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
		return errUnmarshal
	}

	// Delete activity
	var updatedActivities []models.Activity
	for _, activity := range activities {
		if activity.Token != tokenString {
			updatedActivities = append(updatedActivities, activity)
		}
	}

	// Ensure empty array is not null in JSON
	if len(updatedActivities) == 0 {
		updatedActivities = []models.Activity{}
	}

	// Save updated activities
	updatedData, err := json.MarshalIndent(updatedActivities, "", "  ")
	if err != nil {
		log.Println("Error marshalling activities:", err)
		return err
	}
	err = os.WriteFile("data/activities.json", updatedData, 0644)
	if err != nil {
		log.Println("Error writing activities file:", err)
		return err
	}

	return nil
}

func (r *HistoryRepository) FindLoggedByToken(tokenString string) (models.Activity, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Read file
	data, err := os.ReadFile("data/activities.json")
	if err != nil {
		log.Println("Error reading activities file:", err)
	}

	var activities []models.Activity
	errUnmarshal := json.Unmarshal(data, &activities)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
	}

	// find
	for _, activity := range activities {
		if activity.Token == tokenString {
			return activity, true
		}
	}

	return models.Activity{}, false
}

func (r *HistoryRepository) DeleteActivity(activityID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile("data/activities.json")
	if err != nil {
		log.Println("Error reading activities file:", err)
		return err
	}

	var activities []models.Activity
	errUnmarshal := json.Unmarshal(data, &activities)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
		return errUnmarshal
	}

	var updatedActivities []models.Activity
	for _, activity := range activities {
		if activity.ID != activityID {
			updatedActivities = append(updatedActivities, activity)
		}
	}

	// Ensure empty array is not null in JSON
	if len(updatedActivities) == 0 {
		updatedActivities = []models.Activity{}
	}

	updatedData, errMarshal := json.Marshal(updatedActivities)
	if errMarshal != nil {
		log.Fatalf("Error during Marshal: %v", errMarshal)
		return errMarshal
	}

	errWrite := os.WriteFile("data/activities.json", updatedData, 0644)
	if errWrite != nil {
		log.Println("Error writing activities file:", errWrite)
		return errWrite
	}

	return nil
}
