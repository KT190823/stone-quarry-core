package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"mo-da-backend/internal/services"
)

type PrintHandler struct {
	printSvc *services.PrintTemplateService
}

func NewPrintHandler() *PrintHandler {
	return &PrintHandler{printSvc: services.NewPrintTemplateService()}
}

func (h *PrintHandler) List(w http.ResponseWriter, r *http.Request) {
	params := parseListParams(r)
	results, total, err := h.printSvc.List(params)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": total})
}

func (h *PrintHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := h.printSvc.GetByID(id)
	if err != nil {
		http.Error(w, "not found", 404)
		return
	}
	JSON(w, result)
}

func (h *PrintHandler) Create(w http.ResponseWriter, r *http.Request) {
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	if data["id"] == nil || data["id"] == "" {
		data["id"] = newID("PRT")
	}
	if data["code"] == nil || data["code"] == "" {
		if n, ok := data["name"].(string); ok && n != "" {
			data["code"] = "TPL-" + slug(n)
		}
	}
	result, err := h.printSvc.Create(data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(201)
	JSON(w, result)
}

func (h *PrintHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	result, err := h.printSvc.Update(id, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func (h *PrintHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	err := h.printSvc.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *PrintHandler) SetDefault(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// Clear previous default for the same doc type, then set this one.
	ctx := r.Context()
	data, err := readJSON(r)
	if err != nil {
		http.Error(w, "invalid JSON", 400)
		return
	}
	docType, _ := data["docType"].(string)
	if err := h.printSvc.ClearDefault(ctx, docType); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	result, err := h.printSvc.Update(id, map[string]interface{}{"isDefault": true})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	JSON(w, result)
}

func newID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}

func slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var sb strings.Builder
	lastDash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			sb.WriteRune(r)
			lastDash = false
		} else if !lastDash {
			sb.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(sb.String(), "-")
}
