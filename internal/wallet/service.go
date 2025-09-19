package wallet

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sudo-init-do/okies_core/pkg/db"
)

type Service interface {
	GetBalance(ctx context.Context, userID string) (*Wallet, error)
	Fund(ctx context.Context, userID string, amount int64) error
	Withdraw(ctx context.Context, userID string, amount int64) error
}

type service struct {
	repo Repository
}

func NewService() Service {
	return &service{repo: NewRepository(db.DB)}
}

func (s *service) GetBalance(ctx context.Context, userID string) (*Wallet, error) {
	wallet, err := s.repo.GetByUserID(ctx, userID)
	if err == sql.ErrNoRows {
		// Fallback: create wallet if missing
		return s.repo.CreateWallet(ctx, userID)
	}
	return wallet, err
}

func (s *service) Fund(ctx context.Context, userID string, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	return s.repo.UpdateBalance(ctx, userID, amount)
}

func (s *service) Withdraw(ctx context.Context, userID string, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	// check balance first
	wallet, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if wallet.Balance < amount {
		return errors.New("insufficient balance")
	}

	return s.repo.UpdateBalance(ctx, userID, -amount)
}
