package handlers

import (
	"net/http"
	"strings"

	"mo-da-backend/internal/database"
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

	// The API may start while its PostgreSQL container is offline. Without this
	// guard, Authenticate dereferences a nil pool and the connection closes,
	// which the frontend proxy presents as a 504 Gateway Timeout.
	if database.Pool == nil {
		JSONError(w, http.StatusServiceUnavailable,
			"Cơ sở dữ liệu chưa sẵn sàng. Vui lòng khởi động PostgreSQL và thử lại.")
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

	menu, _ := services.MenuForSession(r.Context(), s)

	JSON(w, map[string]interface{}{
		"token":   token,
		"user":    services.SessionToMap(s),
		"expires": s.ExpiresAt,
		"menu":    menu,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	services.DeleteSession(bearerToken(r))
	JSON(w, map[string]interface{}{"status": "ok"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	s, err := services.GetSession(bearerToken(r))
	if err != nil {
		JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	menu, _ := services.MenuForSession(r.Context(), s)
	JSON(w, map[string]interface{}{
		"user": services.SessionToMap(s),
		"menu": menu,
	})
}

func (h *AuthHandler) Menu(w http.ResponseWriter, r *http.Request) {
	s, err := services.GetSession(bearerToken(r))
	if err != nil {
		JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}
	menu, err := services.MenuForSession(r.Context(), s)
	if err != nil {
		JSON(w, map[string]interface{}{"error": err.Error(), "sections": []interface{}{}})
		return
	}
	JSON(w, menu)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	s, err := services.GetSession(bearerToken(r))
	if err != nil {
		JSONError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var body struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := decodeJSON(r, &body); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ")
		return
	}

	if err := services.ChangePassword(r.Context(), s.Username, body.CurrentPassword, body.NewPassword); err != nil {
		JSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	JSON(w, map[string]interface{}{"status": "ok"})
}

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if strings.HasPrefix(h, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(h, "Bearer "))
	}
	return ""
}
