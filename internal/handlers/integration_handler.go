package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"mo-da-backend/internal/models"
	"mo-da-backend/internal/services"
)

type IntegrationHandler struct {
	svc *services.IntegrationService
}

func NewIntegrationHandler() *IntegrationHandler {
	return &IntegrationHandler{svc: services.NewIntegrationService()}
}

// POST /api/v1/integration/webhook
// Ingests real-time payloads from 3rd party mine software, DJI Terra, LiDAR scanners, or weighbridges
func (h *IntegrationHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Không thể đọc nội dung request")
		return
	}

	// Support both standard { source, eventType, externalId, payload } and direct payload
	var req models.WebhookRequest
	if err := json.Unmarshal(body, &req); err != nil {
		JSONError(w, http.StatusBadRequest, "Định dạng JSON không hợp lệ: "+err.Error())
		return
	}

	if len(req.Payload) == 0 {
		// If root object is the payload itself
		req.Payload = body
		if req.Source == "" {
			req.Source = "MineSurveySoftware"
		}
		if req.EventType == "" {
			req.EventType = "survey_completed"
		}
	}

	event, err := h.svc.ProcessWebhook(&req)
	if err != nil {
		JSON(w, map[string]interface{}{
			"success": false,
			"status":  "failed",
			"eventId": event.ID,
			"error":   err.Error(),
		})
		return
	}

	JSON(w, map[string]interface{}{
		"success":     true,
		"status":      "completed",
		"eventId":     event.ID,
		"receivedAt":  event.ReceivedAt,
		"processedAt": event.ProcessedAt,
		"message":     "Đã tiếp nhận và đồng bộ dữ liệu vào hệ thống mỏ đá thành công!",
	})
}

// GET /api/v1/integration/events?limit=30
func (h *IntegrationHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 30
	}
	list, err := h.svc.ListEvents(limit)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, list)
}

// POST /api/v1/integration/sync/mock-drone?quarryCode=MO-PT-01
// Simulates an incoming 3D Drone RTK scan with surface model metadata and volume calculations
func (h *IntegrationHandler) TriggerMockDrone(w http.ResponseWriter, r *http.Request) {
	quarryCode := r.URL.Query().Get("quarryCode")
	if quarryCode == "" {
		quarryCode = "MO-PT-01"
	}
	event, err := h.svc.TriggerMockDroneSync(quarryCode)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Lỗi đồng bộ mô phỏng: "+err.Error())
		return
	}
	JSON(w, map[string]interface{}{
		"success": true,
		"status":  "completed",
		"eventId": event.ID,
		"message": "Đã tạo chu kỳ đo đạc 3D Drone RTK mô phỏng và tự động đối soát sản lượng thành công!",
	})
}

// Register registers all Integration Gateway endpoints onto the router
func (h *IntegrationHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/integration/webhook", h.Webhook)
	mux.HandleFunc("GET /api/v1/integration/events", h.ListEvents)
	mux.HandleFunc("POST /api/v1/integration/sync/mock-drone", h.TriggerMockDrone)
}
