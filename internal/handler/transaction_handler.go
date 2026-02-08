package handler

import (
	"encoding/json"
	"kasir-api/internal/domain"
	"kasir-api/internal/service"
	"net/http"
)

type TransactionHandler struct {
	service *service.TransactionService
}

func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Checkout godoc
// @Summary      Process checkout
// @Description  Process a checkout transaction with multiple items
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        checkout  body      domain.CheckoutRequest  true  "Checkout data"
// @Success      200       {object}  domain.Transaction
// @Failure      400       {string}  string  "Invalid request payload"
// @Failure      500       {string}  string  "Failed to process checkout"
// @Router       /checkout [post]
func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req domain.CheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	transaction, err := h.service.Checkout(req.Items)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "Failed to process checkout")
		return
	}

	WriteJSON(w, http.StatusOK, transaction)
}
