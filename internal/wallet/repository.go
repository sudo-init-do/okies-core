package wallet

import (
	"context"
	"database/sql"
)

type Repository interface {
	CreateWallet(ctx context.Context, userID string) (*Wallet, error)
	GetByUserID(ctx context.Context, userID string) (*Wallet, error)
	UpdateBalance(ctx context.Context, userID string, amount int64) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

// ✅ Create a new wallet with balance = 0
func (r *repository) CreateWallet(ctx context.Context, userID string) (*Wallet, error) {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO wallets (user_id, balance) VALUES ($1, 0)`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		UserID:  userID,
		Balance: 0,
	}, nil
}

// ✅ Get wallet by userID
func (r *repository) GetByUserID(ctx context.Context, userID string) (*Wallet, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT user_id, balance FROM wallets WHERE user_id = $1`,
		userID,
	)

	var w Wallet
	if err := row.Scan(&w.UserID, &w.Balance); err != nil {
		return nil, err
	}
	return &w, nil
}

// ✅ Update wallet balance (add or subtract)
func (r *repository) UpdateBalance(ctx context.Context, userID string, amount int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE wallets SET balance = balance + $1 WHERE user_id = $2`,
		amount, userID,
	)
	return err
}
