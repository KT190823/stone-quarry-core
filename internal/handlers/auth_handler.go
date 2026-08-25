package handlers

import (
	"net/http"
	"strings"

	"mo-da-backend/internal/services"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decodeJSON(r, &body); err != nil {
		JSON(w, map[string]interface{}{"error": "Dữ liệu không hợp lệ"})
		return
	}

	s, err := services.Authenticate(r.Context(), body.Username, body.Password)
	if err != nil {
		JSON(w, map[string]interface{}{"error": err.Error()})
		return
	}

	token, err := services.CreateSession(s)
	if err != nil {
		JSON(w, map[string]interface{}{"error": "Không thể tạo phiên đăng nhập"})
		return
	}

	JSON(w, map[string]interface{}{
		"token":   token,
		"user":    services.SessionToMap(s),
		"expires": s.ExpiresAt,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	services.DeleteSession(bearerToken(r))
	JSON(w, map[string]interface{}{"status": "ok"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	s, err := services.GetSession(bearerToken(r))
	if err != nil {
		JSON(w, map[string]interface{}{"error": err.Error()})
		return
	}
	JSON(w, services.SessionToMap(s))
}

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
	}
	return ""
}
