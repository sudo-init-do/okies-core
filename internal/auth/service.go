package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sudo-init-do/okies_core/pkg/utils"
	"github.com/sudo-init-do/okies_core/internal/wallet"
)

type Service interface {
	Signup(ctx context.Context, req SignupRequest) error
	Login(ctx context.Context, req LoginRequest) (string, error)
}

type service struct {
	repo       Repository
	walletRepo wallet.Repository
}

func NewService(repo Repository, walletRepo wallet.Repository) Service {
	return &service{repo: repo, walletRepo: walletRepo}
}

func (s *service) Signup(ctx context.Context, req SignupRequest) error {
	// Hash password
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hash,
		Role:         "user",
	}

	// Create user
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return err
	}

	// ✅ Auto-create wallet with balance = 0
	_, err = s.walletRepo.CreateWallet(ctx, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (string, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
