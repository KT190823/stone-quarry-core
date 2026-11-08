package handlers

import (
	"net/http"
	"strings"
	"time"

	"mo-da-backend/internal/database"
)

// OrganizationHandler exposes CRUD for the companies/corporations that own
// the quarries (Tổ chức > Công ty).
type OrganizationHandler struct{}

func NewOrganizationHandler() *OrganizationHandler {
	return &OrganizationHandler{}
}

func (h *OrganizationHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/organizations", h.List)
	mux.HandleFunc("GET /api/organizations/{id}", h.Get)
	mux.HandleFunc("POST /api/organizations", h.Create)
	mux.HandleFunc("PUT /api/organizations/{id}", h.Update)
	mux.HandleFunc("DELETE /api/organizations/{id}", h.Delete)
}

type organizationPayload struct {
	Code    string `json:"code"`
	Name    string `json:"name"`
	TaxCode string `json:"taxCode"`
	Address string `json:"address"`
	Status  string `json:"status"`
}

func (h *OrganizationHandler) List(w http.ResponseWriter, r *http.Request) {
	if database.Pool == nil {
		JSONError(w, http.StatusServiceUnavailable, "Chưa kết nối cơ sở dữ liệu")
		return
	}
	search := strings.TrimSpace(r.URL.Query().Get("search"))

	rows, err := database.Pool.Query(r.Context(), `
		SELECT id::text, code, name, COALESCE(tax_code, ''), COALESCE(address, ''),
		       COALESCE(status, 'active'), created_at, updated_at
		FROM organizations
		WHERE ($1 = '' OR code ILIKE '%' || $1 || '%' OR name ILIKE '%' || $1 || '%'
		       OR COALESCE(tax_code, '') ILIKE '%' || $1 || '%')
		ORDER BY name ASC
	`, search)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var id, code, name, taxCode, address, status string
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&id, &code, &name, &taxCode, &address, &status, &createdAt, &updatedAt); err != nil {
			JSONError(w, http.StatusInternalServerError, err.Error())
			return
		}
		results = append(results, organizationMap(id, code, name, taxCode, address, status, createdAt, updatedAt))
	}
	JSON(w, map[string]interface{}{"data": results, "total": len(results)})
}

func (h *OrganizationHandler) Get(w http.ResponseWriter, r *http.Request) {
	if database.Pool == nil {
		JSONError(w, http.StatusServiceUnavailable, "Chưa kết nối cơ sở dữ liệu")
		return
	}
	idOrCode := r.PathValue("id")
	var id, code, name, taxCode, address, status string
	var createdAt, updatedAt time.Time
	err := database.Pool.QueryRow(r.Context(), `
		SELECT id::text, code, name, COALESCE(tax_code, ''), COALESCE(address, ''),
		       COALESCE(status, 'active'), created_at, updated_at
		FROM organizations
		WHERE id::text = $1 OR code = $1
	`, idOrCode).Scan(&id, &code, &name, &taxCode, &address, &status, &createdAt, &updatedAt)
	if err != nil {
		JSONError(w, http.StatusNotFound, "Không tìm thấy công ty")
		return
	}
	JSON(w, map[string]interface{}{"data": organizationMap(id, code, name, taxCode, address, status, createdAt, updatedAt)})
}

func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	if database.Pool == nil {
		JSONError(w, http.StatusServiceUnavailable, "Chưa kết nối cơ sở dữ liệu")
		return
	}
	var body organizationPayload
	if err := decodeJSON(r, &body); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ")
		return
	}
	body.Code = strings.TrimSpace(body.Code)
	body.Name = strings.TrimSpace(body.Name)
	if body.Code == "" || body.Name == "" {
		JSONError(w, http.StatusBadRequest, "Mã và tên công ty là bắt buộc")
		return
	}
	if body.Status == "" {
		body.Status = "active"
	}

	var id, code, name, taxCode, address, status string
	var createdAt, updatedAt time.Time
	err := database.Pool.QueryRow(r.Context(), `
		INSERT INTO organizations (code, name, tax_code, address, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id::text, code, name, COALESCE(tax_code, ''), COALESCE(address, ''),
		          COALESCE(status, 'active'), created_at, updated_at
	`, body.Code, body.Name, body.TaxCode, body.Address, body.Status).
		Scan(&id, &code, &name, &taxCode, &address, &status, &createdAt, &updatedAt)
	if err != nil {
		JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	JSON(w, map[string]interface{}{"data": organizationMap(id, code, name, taxCode, address, status, createdAt, updatedAt)})
}

func (h *OrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
	if database.Pool == nil {
		JSONError(w, http.StatusServiceUnavailable, "Chưa kết nối cơ sở dữ liệu")
		return
	}
	idOrCode := r.PathValue("id")
	var body organizationPayload
	if err := decodeJSON(r, &body); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ")
		return
	}
	body.Code = strings.TrimSpace(body.Code)
	body.Name = strings.TrimSpace(body.Name)
	if body.Code == "" || body.Name == "" {
		JSONError(w, http.StatusBadRequest, "Mã và tên công ty là bắt buộc")
		return
	}
	if body.Status == "" {
		body.Status = "active"
	}

	var id, code, name, taxCode, address, status string
	var createdAt, updatedAt time.Time
	err := database.Pool.QueryRow(r.Context(), `
		UPDATE organizations
		SET code = $2, name = $3, tax_code = $4, address = $5, status = $6, updated_at = NOW()
		WHERE id::text = $1 OR code = $1
		RETURNING id::text, code, name, COALESCE(tax_code, ''), COALESCE(address, ''),
		          COALESCE(status, 'active'), created_at, updated_at
	`, idOrCode, body.Code, body.Name, body.TaxCode, body.Address, body.Status).
		Scan(&id, &code, &name, &taxCode, &address, &status, &createdAt, &updatedAt)
	if err != nil {
		JSONError(w, http.StatusNotFound, "Không tìm thấy công ty để cập nhật")
		return
	}
	JSON(w, map[string]interface{}{"data": organizationMap(id, code, name, taxCode, address, status, createdAt, updatedAt)})
}

func (h *OrganizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if database.Pool == nil {
		JSONError(w, http.StatusServiceUnavailable, "Chưa kết nối cơ sở dữ liệu")
		return
	}
	idOrCode := r.PathValue("id")
	tag, err := database.Pool.Exec(r.Context(), `
		DELETE FROM organizations WHERE id::text = $1 OR code = $1
	`, idOrCode)
	if err != nil {
		JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if tag.RowsAffected() == 0 {
		JSONError(w, http.StatusNotFound, "Không tìm thấy công ty")
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func organizationMap(id, code, name, taxCode, address, status string, createdAt, updatedAt time.Time) map[string]interface{} {
	return map[string]interface{}{
		"id":        id,
		"code":      code,
		"name":      name,
		"taxCode":   taxCode,
		"address":   address,
		"status":    status,
		"createdAt": createdAt,
		"updatedAt": updatedAt,
	}
}
