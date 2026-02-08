package service

import (
	"kasir-api/internal/domain"
	"kasir-api/internal/repository"
)

type TransactionService struct {
	repo repository.TransactionRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(repo repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []domain.CheckoutItem) (*domain.Transaction, error) {
	return s.repo.CreateTransaction(items)
}