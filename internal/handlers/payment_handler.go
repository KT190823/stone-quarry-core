package handlers

import (
	"net/http"
	"time"

	"mo-da-backend/internal/database"
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
	ctx := r.Context()
	rows, err := database.Pool.Query(ctx, `SELECT id, code, COALESCE(partner, '') as partner, COALESCE(customer, '') as customer, COALESCE(amount, '0 đ') as amount, COALESCE(balance, '0 đ') as balance, COALESCE(due, '') as due, COALESCE(status, 'paid') as status, created_at FROM payments_invoices ORDER BY id DESC`)
	if err != nil {
		JSON(w, []interface{}{})
		return
	}
	defer rows.Close()

	var list []map[string]interface{}
	for rows.Next() {
		var id int
		var code, partner, customer, amount, balance, due, status string
		var createdAt time.Time
		if err := rows.Scan(&id, &code, &partner, &customer, &amount, &balance, &due, &status, &createdAt); err == nil {
			list = append(list, map[string]interface{}{
				"id":         id,
				"code":       code,
				"partner":    partner,
				"customer":   customer,
				"amount":     amount,
				"balance":    balance,
				"due":        due,
				"status":     status,
				"created_at": createdAt.Format(time.RFC3339),
			})
		}
	}
	if list == nil {
		list = []map[string]interface{}{}
	}
	JSON(w, list)
}
