package handlers

import (
	"net/http"
	"strconv"

	"mo-da-backend/internal/services"
)

// CollaborationHandler exposes per-detail comments and activity history.
type CollaborationHandler struct {
	svc *services.CollaborationService
}

func NewCollaborationHandler() *CollaborationHandler {
	return &CollaborationHandler{svc: services.NewCollaborationService()}
}

func (h *CollaborationHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/collaboration/comments", h.ListComments)
	mux.HandleFunc("POST /api/collaboration/comments", h.CreateComment)
	mux.HandleFunc("PUT /api/collaboration/comments/{id}", h.UpdateComment)
	mux.HandleFunc("DELETE /api/collaboration/comments/{id}", h.DeleteComment)
	mux.HandleFunc("GET /api/collaboration/activities", h.ListActivities)
	mux.HandleFunc("POST /api/collaboration/activities", h.CreateActivity)
}

type collaborationSession struct {
	Username string
	FullName string
	ClientID string
}

func sessionForRequest(r *http.Request) (collaborationSession, bool) {
	s, err := services.GetSession(bearerToken(r))
	if err != nil {
		return collaborationSession{}, false
	}
	m := services.SessionToMap(s)
	username, _ := m["username"].(string)
	fullName, _ := m["name"].(string)
	clientID, _ := m["clientId"].(string)
	if fullName == "" {
		fullName = username
	}
	return collaborationSession{Username: username, FullName: fullName, ClientID: clientID}, username != ""
}

func (h *CollaborationHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	sess, ok := sessionForRequest(r)
	if !ok {
		JSONError(w, http.StatusUnauthorized, "Phiên đăng nhập không hợp lệ")
		return
	}
	results, err := h.svc.ListComments(r.Context(), sess.ClientID,
		r.URL.Query().Get("entityType"), r.URL.Query().Get("entityId"))
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": len(results)})
}

func (h *CollaborationHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	sess, ok := sessionForRequest(r)
	if !ok {
		JSONError(w, http.StatusUnauthorized, "Phiên đăng nhập không hợp lệ")
		return
	}
	var body struct {
		EntityType string `json:"entityType"`
		EntityID   string `json:"entityId"`
		Content    string `json:"content"`
		ParentID   int    `json:"parentId"`
	}
	if err := decodeJSON(r, &body); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ")
		return
	}
	result, err := h.svc.CreateComment(r.Context(), sess.ClientID, body.EntityType, body.EntityID,
		sess.Username, sess.FullName, body.Content, body.ParentID)
	if err != nil {
		JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	JSON(w, result)
}

func (h *CollaborationHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	sess, ok := sessionForRequest(r)
	if !ok {
		JSONError(w, http.StatusUnauthorized, "Phiên đăng nhập không hợp lệ")
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Bình luận không hợp lệ")
		return
	}
	var body struct {
		Content string `json:"content"`
	}
	if err := decodeJSON(r, &body); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ")
		return
	}
	result, err := h.svc.UpdateComment(r.Context(), sess.ClientID, id, sess.Username, body.Content)
	if err != nil {
		JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON(w, result)
}

func (h *CollaborationHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	sess, ok := sessionForRequest(r)
	if !ok {
		JSONError(w, http.StatusUnauthorized, "Phiên đăng nhập không hợp lệ")
		return
	}
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Bình luận không hợp lệ")
		return
	}
	if err := h.svc.DeleteComment(r.Context(), sess.ClientID, id, sess.Username); err != nil {
		JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON(w, map[string]string{"status": "deleted"})
}

func (h *CollaborationHandler) ListActivities(w http.ResponseWriter, r *http.Request) {
	sess, ok := sessionForRequest(r)
	if !ok {
		JSONError(w, http.StatusUnauthorized, "Phiên đăng nhập không hợp lệ")
		return
	}
	results, err := h.svc.ListActivities(r.Context(), sess.ClientID,
		r.URL.Query().Get("entityType"), r.URL.Query().Get("entityId"))
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	JSON(w, map[string]interface{}{"data": results, "total": len(results)})
}

func (h *CollaborationHandler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	sess, ok := sessionForRequest(r)
	if !ok {
		JSONError(w, http.StatusUnauthorized, "Phiên đăng nhập không hợp lệ")
		return
	}
	var body struct {
		EntityType string `json:"entityType"`
		EntityID   string `json:"entityId"`
		Action     string `json:"action"`
		ActionType string `json:"actionType"`
		Detail     string `json:"detail"`
	}
	if err := decodeJSON(r, &body); err != nil {
		JSONError(w, http.StatusBadRequest, "Dữ liệu không hợp lệ")
		return
	}
	result, err := h.svc.CreateActivity(r.Context(), sess.ClientID, body.EntityType, body.EntityID,
		sess.Username, sess.FullName, body.Action, body.ActionType, body.Detail)
	if err != nil {
		JSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	JSON(w, result)
}
