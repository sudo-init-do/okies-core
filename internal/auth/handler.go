package auth

import (
	"encoding/json"
	"net/http"

	"github.com/sudo-init-do/okies_core/pkg/db"
	"github.com/sudo-init-do/okies_core/pkg/response"
	"github.com/sudo-init-do/okies_core/internal/wallet"
)

type Handler struct {
	service Service
}

// Now wires both user repo and wallet repo
func NewHandler() *Handler {
	repo := NewRepository(db.DB)
	walletRepo := wallet.NewRepository(db.DB)
	svc := NewService(repo, walletRepo)
	return &Handler{service: svc}
}

// Signup endpoint
func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Write(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := h.service.Signup(r.Context(), req); err != nil {
		response.Write(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Write(w, http.StatusCreated, "user created successfully", nil)
}

// Login endpoint
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Write(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	token, err := h.service.Login(r.Context(), req)
	if err != nil {
		response.Write(w, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Write(w, http.StatusOK, "login successful", map[string]string{
		"token": token,
	})
}
