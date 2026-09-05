package handlers

import (
	"net/http"

	"mo-da-backend/internal/services"
)

type CameraHandler struct {
	cameraSvc     *services.CameraDeviceService
	googleDrive   *services.GoogleDriveService
	googleEmail   *services.GoogleEmailService
	googleMap     *services.GoogleMapService
	googlePhoto   *services.GooglePhotoService
	gpsFleet      *services.GpsFleetService
	fuelTheft     *services.FuelTheftAuditService
	fuelNorm      *services.FuelNormConfigService
	yardCheck     *services.YardCheckService
}

func NewCameraHandler() *CameraHandler {
	return &CameraHandler{
		cameraSvc:   services.NewCameraDeviceService(),
		googleDrive: services.NewGoogleDriveService(),
		googleEmail: services.NewGoogleEmailService(),
		googleMap:   services.NewGoogleMapService(),
		googlePhoto: services.NewGooglePhotoService(),
		gpsFleet:    services.NewGpsFleetService(),
		fuelTheft:   services.NewFuelTheftAuditService(),
		fuelNorm:    services.NewFuelNormConfigService(),
		yardCheck:   services.NewYardCheckService(),
	}
}

func (h *CameraHandler) ListCameras(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.cameraSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CameraHandler) GetCamera(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.cameraSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *CameraHandler) CreateCamera(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.cameraSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *CameraHandler) UpdateCamera(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.cameraSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *CameraHandler) DeleteCamera(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.cameraSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *CameraHandler) ListGpsFleet(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.gpsFleet.ListJSONB(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CameraHandler) ListFuelTheftAudits(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.fuelTheft.ListJSONB(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CameraHandler) GetFuelTheftAudit(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.fuelTheft.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}


func (h *CameraHandler) ListFuelNorms(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.fuelNorm.ListJSONB(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CameraHandler) ListYardCheckInOuts(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.yardCheck.ListJSONB(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CameraHandler) ListGoogleDrive(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.googleDrive.ListJSONB(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CameraHandler) ListGoogleGmail(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.googleEmail.ListJSONB(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CameraHandler) ListGoogleMaps(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.googleMap.ListJSONB(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *CameraHandler) ListGooglePhotos(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.googlePhoto.ListJSONB(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}
