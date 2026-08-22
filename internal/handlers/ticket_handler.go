package handlers

import (
	"net/http"

	"mo-da-backend/internal/database"
	"mo-da-backend/internal/services"
)

type TicketHandler struct {
	svc *services.TicketService
}

func NewTicketHandler() *TicketHandler {
	return &TicketHandler{svc: services.NewTicketService()}
}

func (h *TicketHandler) List(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.svc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *TicketHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.svc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.svc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *TicketHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.svc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *TicketHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.svc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	var ticketCount, vehicleCount, employeeCount int
	database.Pool.QueryRow(r.Context(), "SELECT COUNT(*) FROM tickets").Scan(&ticketCount)
	database.Pool.QueryRow(r.Context(), "SELECT COUNT(*) FROM vehicles").Scan(&vehicleCount)
	database.Pool.QueryRow(r.Context(), "SELECT COUNT(*) FROM hr_employees").Scan(&employeeCount)

	JSON(w, map[string]interface{}{
		"stats": map[string]interface{}{
			"tickets":   ticketCount,
			"vehicles":  vehicleCount,
			"employees": employeeCount,
		},
	})
}
