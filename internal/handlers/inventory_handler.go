package handlers

import (
	"net/http"

	"mo-da-backend/internal/services"
)

type InventoryHandler struct {
	inboundSvc    *services.InventoryInboundService
	outboundSvc   *services.InventoryOutboundService
	stocktakeSvc  *services.InventoryStocktakeService
	movementSvc   *services.InventoryMovementService
}

func NewInventoryHandler() *InventoryHandler {
	return &InventoryHandler{
		inboundSvc:   services.NewInventoryInboundService(),
		outboundSvc:  services.NewInventoryOutboundService(),
		stocktakeSvc: services.NewInventoryStocktakeService(),
		movementSvc:  services.NewInventoryMovementService(),
	}
}

func (h *InventoryHandler) ListInbound(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.inboundSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *InventoryHandler) ListOutbound(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.outboundSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *InventoryHandler) ListStocktake(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.stocktakeSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *InventoryHandler) ListMovements(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.movementSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}
