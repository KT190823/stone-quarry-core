package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"mo-da-backend/internal/models"
	"mo-da-backend/internal/services"
)

type QuarryHandler struct {
	svc *services.QuarryService
}

func NewQuarryHandler() *QuarryHandler {
	return &QuarryHandler{svc: services.NewQuarryService()}
}

// GET /api/v1/quarries
func (h *QuarryHandler) ListQuarries(w http.ResponseWriter, r *http.Request) {
	list, err := h.svc.ListQuarries()
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, list)
}

// GET /api/v1/quarries/{id}
func (h *QuarryHandler) GetQuarry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	q, err := h.svc.GetQuarryByID(id)
	if err != nil {
		JSONError(w, http.StatusNotFound, "Không tìm thấy mỏ đá")
		return
	}
	JSON(w, q)
}

// POST /api/v1/quarries
func (h *QuarryHandler) CreateQuarry(w http.ResponseWriter, r *http.Request) {
	var q models.Quarry
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ")
		return
	}
	if err := h.svc.CreateQuarry(&q); err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, q)
}

// GET /api/v1/quarry-areas?quarryId=...
func (h *QuarryHandler) ListQuarryAreas(w http.ResponseWriter, r *http.Request) {
	quarryID := r.URL.Query().Get("quarryId")
	list, err := h.svc.ListQuarryAreas(quarryID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, list)
}

// POST /api/v1/quarry-areas
func (h *QuarryHandler) CreateQuarryArea(w http.ResponseWriter, r *http.Request) {
	var a models.QuarryArea
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ")
		return
	}
	if err := h.svc.CreateQuarryArea(&a); err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, a)
}

// GET /api/v1/survey-cycles?quarryId=...
func (h *QuarryHandler) ListSurveyCycles(w http.ResponseWriter, r *http.Request) {
	quarryID := r.URL.Query().Get("quarryId")
	list, err := h.svc.ListSurveyCycles(quarryID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, list)
}

// GET /api/v1/survey-cycles/{id}
func (h *QuarryHandler) GetSurveyCycle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cycle, err := h.svc.GetSurveyCycleByID(id)
	if err != nil {
		JSONError(w, http.StatusNotFound, "Không tìm thấy chu kỳ đo đạc")
		return
	}
	JSON(w, cycle)
}

// POST /api/v1/survey-cycles
func (h *QuarryHandler) CreateSurveyCycle(w http.ResponseWriter, r *http.Request) {
	var cycle models.SurveyCycle
	if err := json.NewDecoder(r.Body).Decode(&cycle); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ")
		return
	}
	if err := h.svc.CreateSurveyCycle(&cycle); err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, cycle)
}

// GET /api/v1/volume-calculations?cycleId=...
func (h *QuarryHandler) ListVolumeCalculations(w http.ResponseWriter, r *http.Request) {
	cycleID := r.URL.Query().Get("cycleId")
	list, err := h.svc.ListVolumeCalculations(cycleID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, list)
}

// GET /api/v1/surface-models?cycleId=...
func (h *QuarryHandler) ListSurfaceModels(w http.ResponseWriter, r *http.Request) {
	cycleID := r.URL.Query().Get("cycleId")
	list, err := h.svc.ListSurfaceModels(cycleID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, list)
}

// GET /api/v1/quarry-materials?quarryId=...
func (h *QuarryHandler) ListMaterials(w http.ResponseWriter, r *http.Request) {
	quarryID := r.URL.Query().Get("quarryId")
	list, err := h.svc.ListMaterials(quarryID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, list)
}

// GET /api/v1/weighing-transactions?quarryId=...&cycleId=...&page=1&pageSize=50
func (h *QuarryHandler) ListWeighingTransactions(w http.ResponseWriter, r *http.Request) {
	quarryID := r.URL.Query().Get("quarryId")
	cycleID := r.URL.Query().Get("cycleId")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	offset := (page - 1) * pageSize

	list, total, err := h.svc.ListWeighingTransactions(quarryID, cycleID, pageSize, offset)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, map[string]interface{}{
		"data":     list,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// GET /api/v1/reconciliations?quarryId=...
func (h *QuarryHandler) ListReconciliations(w http.ResponseWriter, r *http.Request) {
	quarryID := r.URL.Query().Get("quarryId")
	list, err := h.svc.ListReconciliations(quarryID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, list)
}

// GET /api/v1/quarry-alerts?quarryId=...
func (h *QuarryHandler) ListAlerts(w http.ResponseWriter, r *http.Request) {
	quarryID := r.URL.Query().Get("quarryId")
	list, err := h.svc.ListAlerts(quarryID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, list)
}

// GET /api/v1/quarry/dashboard?quarryId=...
func (h *QuarryHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	quarryID := r.URL.Query().Get("quarryId")
	summary, err := h.svc.GetDashboardSummary(quarryID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, summary)
}

// Register registers all Quarry REST endpoints onto the HTTP router
func (h *QuarryHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/quarries", h.ListQuarries)
	mux.HandleFunc("GET /api/v1/quarries/{id}", h.GetQuarry)
	mux.HandleFunc("POST /api/v1/quarries", h.CreateQuarry)

	mux.HandleFunc("GET /api/v1/quarry-areas", h.ListQuarryAreas)
	mux.HandleFunc("POST /api/v1/quarry-areas", h.CreateQuarryArea)

	mux.HandleFunc("GET /api/v1/survey-cycles", h.ListSurveyCycles)
	mux.HandleFunc("GET /api/v1/survey-cycles/{id}", h.GetSurveyCycle)
	mux.HandleFunc("POST /api/v1/survey-cycles", h.CreateSurveyCycle)

	mux.HandleFunc("GET /api/v1/volume-calculations", h.ListVolumeCalculations)
	mux.HandleFunc("GET /api/v1/surface-models", h.ListSurfaceModels)
	mux.HandleFunc("GET /api/v1/quarry-materials", h.ListMaterials)

	mux.HandleFunc("GET /api/v1/weighing-transactions", h.ListWeighingTransactions)
	mux.HandleFunc("GET /api/v1/reconciliations", h.ListReconciliations)
	mux.HandleFunc("GET /api/v1/quarry-alerts", h.ListAlerts)

	mux.HandleFunc("GET /api/v1/quarry/dashboard", h.Dashboard)
}

// Helper to sanitize path params if needed
func trimSlash(s string) string {
	return strings.Trim(s, "/")
}
