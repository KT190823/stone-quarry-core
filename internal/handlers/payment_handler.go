package handlers

import (
	"net/http"

	"mo-da-backend/internal/services"
)

type PaymentHandler struct {
	debtSvc      *services.DebtService
	reconcileSvc *services.ReconcileService
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		debtSvc:      services.NewDebtService(),
		reconcileSvc: services.NewReconcileService(),
	}
}

func (h *PaymentHandler) ListDebt(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.debtSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *PaymentHandler) ListReconcile(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.reconcileSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *PaymentHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	JSON(w, []interface{}{})
}
