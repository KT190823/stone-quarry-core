package handlers

import (
	"net/http"

	"mo-da-backend/internal/services"
)

type UserHandler struct {
	userSvc     *services.UserService
	roleSvc     *services.UserRoleService
	logSvc      *services.UserLogService
	reportSvc   *services.ReportService
	settingSvc  *services.SettingService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userSvc:    services.NewUserService(),
		roleSvc:    services.NewUserRoleService(),
		logSvc:     services.NewUserLogService(),
		reportSvc:  services.NewReportService(),
		settingSvc: services.NewSettingService(),
	}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.userSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.userSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.userSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.userSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.userSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *UserHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.roleSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *UserHandler) ListLogs(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.logSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *UserHandler) ListReports(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.reportSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *UserHandler) GetReport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.reportSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *UserHandler) CreateReport(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.reportSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *UserHandler) UpdateReport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.reportSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *UserHandler) DeleteReport(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.reportSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *UserHandler) ListSettings(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.settingSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *UserHandler) GetSetting(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.settingSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *UserHandler) CreateSetting(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.settingSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *UserHandler) UpdateSetting(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.settingSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *UserHandler) DeleteSetting(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.settingSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}
