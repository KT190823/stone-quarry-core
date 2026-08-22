package handlers

import (
	"net/http"

	"mo-da-backend/internal/services"
)

type MiningHandler struct {
	permitSvc       *services.MiningPermitService
	planSvc         *services.MiningPlanService
	blastingSvc     *services.BlastingPassportService
	crusherSvc      *services.CrusherPlantService
	fuelLogSvc      *services.EquipmentFuelLogService
	statutorySvc    *services.StatutoryReportService
	resourceTaxSvc  *services.ResourceTaxService
	productionSvc   *services.ProductionStageService
}

func NewMiningHandler() *MiningHandler {
	return &MiningHandler{
		permitSvc:      services.NewMiningPermitService(),
		planSvc:        services.NewMiningPlanService(),
		blastingSvc:    services.NewBlastingPassportService(),
		crusherSvc:     services.NewCrusherPlantService(),
		fuelLogSvc:     services.NewEquipmentFuelLogService(),
		statutorySvc:   services.NewStatutoryReportService(),
		resourceTaxSvc: services.NewResourceTaxService(),
		productionSvc:  services.NewProductionStageService(),
	}
}

func (h *MiningHandler) ListPermits(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.permitSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *MiningHandler) GetPermit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.permitSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *MiningHandler) CreatePermit(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.permitSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *MiningHandler) UpdatePermit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.permitSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *MiningHandler) DeletePermit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.permitSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *MiningHandler) ListPlans(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.planSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *MiningHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.planSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *MiningHandler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.planSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *MiningHandler) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.planSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *MiningHandler) DeletePlan(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.planSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *MiningHandler) ListBlasting(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.blastingSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *MiningHandler) GetBlasting(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.blastingSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *MiningHandler) CreateBlasting(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.blastingSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *MiningHandler) UpdateBlasting(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.blastingSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *MiningHandler) DeleteBlasting(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.blastingSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *MiningHandler) ListCrusherPlants(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.crusherSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *MiningHandler) ListEquipmentFuel(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.fuelLogSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *MiningHandler) GetEquipmentFuel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.fuelLogSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *MiningHandler) CreateEquipmentFuel(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.fuelLogSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *MiningHandler) UpdateEquipmentFuel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.fuelLogSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *MiningHandler) DeleteEquipmentFuel(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.fuelLogSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *MiningHandler) ListStatutoryReports(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.statutorySvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *MiningHandler) ListResourceTaxes(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.resourceTaxSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *MiningHandler) ListProductionStages(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.productionSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}
