package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"mo-da-backend/internal/database"
)

func TestLoginReturnsServiceUnavailableWhenDatabaseIsOffline(t *testing.T) {
	if database.Pool != nil {
		t.Skip("test expects no initialized database pool")
	}

	payload := map[string]string{
		"username": "__health_probe_nonexistent__",
		"password": "not-a-password",
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	NewAuthHandler().Login(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected HTTP 503 while DB is unavailable, got %d", response.Code)
	}
	if !bytes.Contains(response.Body.Bytes(), []byte("error")) {
		t.Fatalf("expected a JSON error response, got %q", response.Body.String())
	}
}
