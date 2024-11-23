package repository

import (
	"encoding/json"
	"log"
	"merchant-bank-api/app/models"
	"os"
	"sync"
)

type JSONRepository struct {
	mu sync.RWMutex
}

func NewJSONRepository() *JSONRepository {
	return &JSONRepository{}
}

func (r *JSONRepository) FindCustomerByUsername(username string) (models.Customer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile("data/customers.json")
	if err != nil {
		log.Println("Error reading customers file:", err)
		return models.Customer{}, false
	}

	var customers []models.Customer
	errUnmarshal := json.Unmarshal(data, &customers)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", err)
		return models.Customer{}, false
	}

	for _, customer := range customers {
		if customer.Username == username {
			return customer, true
		}
	}
	return models.Customer{}, false
}

func (r *JSONRepository) FindCustomerByEmail(email string) (models.Customer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile("data/customers.json")
	if err != nil {
		log.Println("Error reading customers file:", err)
		return models.Customer{}, false
	}

	var customers []models.Customer
	errUnmarshal := json.Unmarshal(data, &customers)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", err)
		return models.Customer{}, false
	}

	for _, customer := range customers {
		if customer.Email == email {
			return customer, true
		}
	}
	return models.Customer{}, false
}

func (r *JSONRepository) SaveActivityHistory(activity models.Activity) {
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

func (r *JSONRepository) SavePaymentHistory(transaction models.Transaction) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Read existing transactions
	data, err := os.ReadFile("data/transactions.json")
	if err != nil {
		log.Println("Error reading transactions file:", err)
		return
	}

	var transactions []models.Transaction
	errUnmarshal := json.Unmarshal(data, &transactions)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
	}

	// Append new transaction
	transactions = append(transactions, transaction)

	// Write back to file
	updatedData, _ := json.MarshalIndent(transactions, "", "  ")
	errWrite := os.WriteFile("data/transactions.json", updatedData, 0644)
	if errWrite != nil {
		log.Fatalf("Error during Unmarshal: %v", errWrite)
	}
}

func (r *JSONRepository) SaveCustomer(customer *models.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Read existing customers
	data, err := os.ReadFile("data/customers.json")
	if err != nil {
		log.Println("Error reading customers file:", err)
		return err
	}
	var customers []models.Customer
	errUnmarshal := json.Unmarshal(data, &customers)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
	}
	// Append new customer
	customers = append(customers, *customer)
	// Write back to file
	updatedData, err := json.MarshalIndent(customers, "", " ")
	if err != nil {
		log.Println("Error marshalling customers:", err)
		return err
	}
	err = os.WriteFile("data/customers.json", updatedData, 0644)
	if err != nil {
		log.Println("Error writing customers file:", err)
		return err
	}
	return nil
}

func (r *JSONRepository) FindCustomerById(id string) (models.Customer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile("data/customers.json")
	if err != nil {
		log.Println("Error reading customers file:", err)
		return models.Customer{}, false
	}

	var customers []models.Customer
	errUnmarshal := json.Unmarshal(data, &customers)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
		return models.Customer{}, false
	}

	for _, customer := range customers {
		if customer.ID == id {
			return customer, true
		}
	}
	return models.Customer{}, false
}

func (r *JSONRepository) FindAllCustomers() []models.Customer {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile("data/customers.json")
	if err != nil {
		log.Println("Error reading customers file:", err)
		return []models.Customer{}
	}

	var customers []models.Customer
	errUnmarshal := json.Unmarshal(data, &customers)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
	}

	return customers
}

func (r *JSONRepository) FindPaymentById(id string) (models.Transaction, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile("data/transactions.json")
	if err != nil {
		log.Println("Error reading transactions file:", err)
		return models.Transaction{}, false
	}

	var transactions []models.Transaction
	errUnmarshal := json.Unmarshal(data, &transactions)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
		return models.Transaction{}, false
	}

	for _, transaction := range transactions {
		if transaction.IDTransaction == id {
			return transaction, true
		}
	}
	return models.Transaction{}, false
}

func (r *JSONRepository) FindAllPayment() []models.Transaction {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := os.ReadFile("data/transactions.json")
	if err != nil {
		log.Println("Error reading transactions file:", err)
		return []models.Transaction{}
	}

	var transactions []models.Transaction
	errUnmarshal := json.Unmarshal(data, &transactions)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
		return []models.Transaction{}
	}

	return transactions

}

func (r *JSONRepository) FindAllActivities() []models.Activity {
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

func (r *JSONRepository) FindActivitiesByUsername(username string) []models.Activity {
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

func (r *JSONRepository) UpdateCustomerBalance(customer models.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := os.ReadFile("data/customers.json")
	if err != nil {
		log.Println("Error reading customers file:", err)
		return err
	}

	var customers []models.Customer
	errUnmarshal := json.Unmarshal(data, &customers)
	if errUnmarshal != nil {
		log.Fatalf("Error during Unmarshal: %v", errUnmarshal)
		return errUnmarshal
	}

	// Update customer balance
	for i, c := range customers {
		if c.ID == customer.ID {
			customers[i].Balance = customer.Balance
			break
		}
	}

	// Write again
	updatedData, errMarshall := json.MarshalIndent(customers, "", "  ")
	if errMarshall != nil {
		log.Println("Error marshalling customers:", errMarshall)
		return errMarshall
	}
	err = os.WriteFile("data/customers.json", updatedData, 0644)
	if err != nil {
		log.Println("Error writing customers file:", err)
		return err
	}
	return nil
}

func (r *JSONRepository) DeleteLoggedByToken(tokenString string) error {
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

func (r *JSONRepository) FindLoggedByToken(tokenString string) (models.Activity, bool) {
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
