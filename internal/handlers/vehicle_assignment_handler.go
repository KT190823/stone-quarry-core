package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"mo-da-backend/internal/services"
)

type VehicleAssignmentHandler struct {
	svc *services.VehicleAssignmentService
}

func NewVehicleAssignmentHandler() *VehicleAssignmentHandler {
	return &VehicleAssignmentHandler{svc: services.NewVehicleAssignmentService()}
}

func (h *VehicleAssignmentHandler) List(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.svc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *VehicleAssignmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.svc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *VehicleAssignmentHandler) GetActiveByPlate(w http.ResponseWriter, r *http.Request) {
	plate := strings.TrimSpace(r.URL.Query().Get("plate"))
	if plate == "" {
		http.Error(w, "missing plate parameter", 400)
		return
	}
	result, err := h.svc.GetActiveByPlate(plate)
	if err != nil {
		JSON(w, map[string]interface{}{"found": false, "plate": plate})
		return
	}
	JSON(w, map[string]interface{}{"found": true, "data": result})
}

func (h *VehicleAssignmentHandler) Create(w http.ResponseWriter, r *http.Request) {
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

func (h *VehicleAssignmentHandler) Update(w http.ResponseWriter, r *http.Request) {
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

func (h *VehicleAssignmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.svc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"success": true})
}

func (h *VehicleAssignmentHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var body struct {
		StartOdo        float64 `json:"startOdo"`
		StartFuelLiters float64 `json:"startFuelLiters"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	result, err := h.svc.CheckIn(r.Context(), id, body.StartOdo, body.StartFuelLiters)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *VehicleAssignmentHandler) CheckOut(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var body struct {
		EndOdo        float64 `json:"endOdo"`
		EndFuelLiters float64 `json:"endFuelLiters"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	result, err := h.svc.CheckOut(r.Context(), id, body.EndOdo, body.EndFuelLiters)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}
