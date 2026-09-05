package handlers

import (
	"encoding/json"
	"net/http"

	"mo-da-backend/internal/models"
	"mo-da-backend/internal/services"
)

type SalesContractHandler struct {
	contractSvc     *services.SalesContractService
	consolidatedSvc *services.ConsolidatedDeliveryOrderService
}

func NewSalesContractHandler() *SalesContractHandler {
	return &SalesContractHandler{
		contractSvc:     services.NewSalesContractService(),
		consolidatedSvc: services.NewConsolidatedDeliveryOrderService(),
	}
}

// ListContracts lists commercial sales contracts
func (h *SalesContractHandler) ListContracts(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.contractSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

// GetContract gets a sales contract by ID or Code
func (h *SalesContractHandler) GetContract(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.contractSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

// CreateContract creates a new sales contract
func (h *SalesContractHandler) CreateContract(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.contractSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

// UpdateContract updates a sales contract
func (h *SalesContractHandler) UpdateContract(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.contractSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

// DeleteContract deletes a sales contract
func (h *SalesContractHandler) DeleteContract(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.contractSvc.Delete(id); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"success": true})
}

// SettleDailyPrepaid triggers the automatic end-of-day batch settlement for customer prepaid wallet
func (h *SalesContractHandler) SettleDailyPrepaid(w http.ResponseWriter, r *http.Request) {
	var body struct {
		CustomerCode string `json:"customerCode"`
		Date         string `json:"date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}

	result, err := h.contractSvc.SettleDailyPrepaid(r.Context(), body.CustomerCode, body.Date)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

// ListConsolidatedOrders lists consolidated dispatch orders
func (h *SalesContractHandler) ListConsolidatedOrders(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.consolidatedSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

// GetConsolidatedOrder gets a consolidated delivery order by ID or Code
func (h *SalesContractHandler) GetConsolidatedOrder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.consolidatedSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

// ConsolidateDeliveryOrders handles grouping multiple scale tickets into an official delivery slip + e-Invoice
func (h *SalesContractHandler) ConsolidateDeliveryOrders(w http.ResponseWriter, r *http.Request) {
	var req models.ConsolidatedDeliveryOrder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}

	result, err := h.consolidatedSvc.ConsolidateDeliveryOrders(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}
