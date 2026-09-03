package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mo-da-backend/internal/services"
)

func JSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func JSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   msg,
		"message": msg,
		"status":  status,
	})
}


func parseListParams(r *http.Request) services.ListParams {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}
	return services.ListParams{
		Page:     page,
		PageSize: pageSize,
		Search:   r.URL.Query().Get("search"),
		Sort:     r.URL.Query().Get("sort"),
		Order:    r.URL.Query().Get("order"),
	}
}

func readJSON(r *http.Request) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	return data, err
}

// decodeJSON decodes the request body into an arbitrary struct target.
func decodeJSON(r *http.Request, target interface{}) error {
	return json.NewDecoder(r.Body).Decode(target)
}
