package repository

import (
	"encoding/json"
	"log"
	"merchant-bank-api/app/models"
	"os"
	"sync"
)

type TransactionRepository struct {
	mu sync.RWMutex
}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

func (r *TransactionRepository) FindPaymentById(id string) (models.Transaction, bool) {
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

func (r *TransactionRepository) FindAllPayment() []models.Transaction {
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

func (r *TransactionRepository) SavePaymentHistory(transaction models.Transaction) {
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
