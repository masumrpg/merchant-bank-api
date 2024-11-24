package repository

import (
	"encoding/json"
	"log"
	"merchant-bank-api/app/models"
	"os"
	"sync"
)

type UserRepository struct {
	mu sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) FindUserByUsername(username string) (models.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile("data/users.json")
	if err != nil {
		log.Println("Error reading users file:", err)
		return models.User{}, false
	}

	var users []models.User
	errUnmarshal := json.Unmarshal(data, &users)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", err)
		return models.User{}, false
	}

	for _, user := range users {
		if user.Username == username {
			return user, true
		}
	}
	return models.User{}, false
}

func (r *UserRepository) FindUserByEmail(email string) (models.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile("data/users.json")
	if err != nil {
		log.Println("Error reading users file:", err)
		return models.User{}, false
	}

	var users []models.User
	errUnmarshal := json.Unmarshal(data, &users)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", err)
		return models.User{}, false
	}

	for _, user := range users {
		if user.Email == email {
			return user, true
		}
	}
	return models.User{}, false
}

func (r *UserRepository) SaveUser(user *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Read existing users
	data, err := os.ReadFile("data/users.json")
	if err != nil {
		log.Println("Error reading users file:", err)
		return err
	}
	var users []models.User
	errUnmarshal := json.Unmarshal(data, &users)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
	}
	// Append new user
	users = append(users, *user)
	// Write back to file
	updatedData, err := json.MarshalIndent(users, "", " ")
	if err != nil {
		log.Println("Error marshalling users:", err)
		return err
	}
	err = os.WriteFile("data/users.json", updatedData, 0644)
	if err != nil {
		log.Println("Error writing users file:", err)
		return err
	}
	return nil
}

func (r *UserRepository) FindUserById(id string) (models.User, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile("data/users.json")
	if err != nil {
		log.Println("Error reading users file:", err)
		return models.User{}, false
	}

	var users []models.User
	errUnmarshal := json.Unmarshal(data, &users)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
		return models.User{}, false
	}

	for _, user := range users {
		if user.ID == id {
			return user, true
		}
	}
	return models.User{}, false
}

func (r *UserRepository) FindAllUser() []models.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile("data/users.json")
	if err != nil {
		log.Println("Error reading users file:", err)
		return []models.User{}
	}

	var users []models.User
	errUnmarshal := json.Unmarshal(data, &users)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
	}

	return users
}

func (r *UserRepository) UpdateUserBalance(user models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile("data/users.json")
	if err != nil {
		log.Println("Error reading users file:", err)
		return err
	}

	var users []models.User
	errUnmarshal := json.Unmarshal(data, &users)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
		return errUnmarshal
	}

	// Update user balance
	for i, c := range users {
		if c.ID == user.ID {
			users[i].Balance = user.Balance
			break
		}
	}

	// Write again
	updatedData, errMarshall := json.MarshalIndent(users, "", "  ")
	if errMarshall != nil {
		log.Println("Error marshalling users:", errMarshall)
		return errMarshall
	}
	err = os.WriteFile("data/users.json", updatedData, 0644)
	if err != nil {
		log.Println("Error writing users file:", err)
		return err
	}
	return nil
}
