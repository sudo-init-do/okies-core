package wallet

import (
	"encoding/json"
	"net/http"

	"github.com/sudo-init-do/okies_core/middleware"
	"github.com/sudo-init-do/okies_core/pkg/response"
)

type Handler struct {
	service Service
}

func NewHandler() *Handler {
	return &Handler{service: NewService()}
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(string)

	wallet, err := h.service.GetBalance(r.Context(), userID)
	if err != nil {
		response.Write(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Write(w, http.StatusOK, "wallet balance fetched", wallet)
}

func (h *Handler) Fund(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(string)

	var req FundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Write(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := h.service.Fund(r.Context(), userID, req.Amount); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Write(w, http.StatusOK, "wallet funded successfully", nil)
}

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.ContextUserID).(string)

	var req WithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Write(w, http.StatusBadRequest, "invalid request body", nil)
		return
	}

	if err := h.service.Withdraw(r.Context(), userID, req.Amount); err != nil {
		response.Write(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Write(w, http.StatusOK, "withdrawal successful", nil)
}
