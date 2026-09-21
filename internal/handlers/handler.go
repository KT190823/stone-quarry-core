package handlers

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"mo-da-backend/internal/services"
)

func JSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if v == nil {
		json.NewEncoder(w).Encode(false)
		return
	}

	val := reflect.ValueOf(v)
	// Unwrap pointers and interfaces
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		if val.IsNil() {
			json.NewEncoder(w).Encode(false)
			return
		}
		val = val.Elem()
	}

	// If top-level value is a slice
	if val.Kind() == reflect.Slice {
		if val.IsNil() {
			json.NewEncoder(w).Encode([]interface{}{})
			return
		}
	}

	// If map contains "data" field
	if m, ok := v.(map[string]interface{}); ok {
		if dataVal, exists := m["data"]; exists {
			if dataVal == nil {
				// If response has "total", it's a paginated list -> return []
				if _, hasTotal := m["total"]; hasTotal {
					m["data"] = []interface{}{}
				} else {
					m["data"] = false
				}
			} else {
				dVal := reflect.ValueOf(dataVal)
				for dVal.Kind() == reflect.Ptr || dVal.Kind() == reflect.Interface {
					if dVal.IsNil() {
						break
					}
					dVal = dVal.Elem()
				}
				if !dVal.IsValid() || ((dVal.Kind() == reflect.Ptr || dVal.Kind() == reflect.Interface) && dVal.IsNil()) {
					if _, hasTotal := m["total"]; hasTotal {
						m["data"] = []interface{}{}
					} else {
						m["data"] = false
					}
				} else if dVal.Kind() == reflect.Slice && dVal.IsNil() {
					m["data"] = []interface{}{}
				}
			}
		}
	}

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
	qCode := r.URL.Query().Get("quarryCode")
	if qCode == "" {
		qCode = r.URL.Query().Get("quarry_code")
	}
	return services.ListParams{
		Page:       page,
		PageSize:   pageSize,
		Search:     r.URL.Query().Get("search"),
		Sort:       r.URL.Query().Get("sort"),
		Order:      r.URL.Query().Get("order"),
		QuarryCode: qCode,
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
