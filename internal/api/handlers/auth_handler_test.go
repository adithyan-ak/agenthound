package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adithyan-ak/agenthound/internal/model"
)

func TestHandleLogin_MissingFields(t *testing.T) {
	h := NewAuthHandler(nil, nil, "secret", nil)
	w := httptest.NewRecorder()
	r := newTestRequest(http.MethodPost, "/api/v1/auth/login", []byte(`{}`))
	h.HandleLogin(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected VALIDATION_ERROR, got %s", resp.Error.Code)
	}
}

func TestHandleCreateUser_InvalidRole(t *testing.T) {
	h := NewAuthHandler(nil, nil, "secret", nil)
	w := httptest.NewRecorder()
	r := newTestRequest(http.MethodPost, "/api/v1/auth/users",
		[]byte(`{"username":"a","password":"b","role":"superadmin"}`))
	h.HandleCreateUser(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected VALIDATION_ERROR, got %s", resp.Error.Code)
	}
}

func TestHandleCreateToken_MissingName(t *testing.T) {
	h := NewAuthHandler(nil, nil, "secret", nil)
	w := httptest.NewRecorder()
	r := newTestRequest(http.MethodPost, "/api/v1/auth/tokens", []byte(`{}`))
	r = withAuthUser(r, &model.User{ID: "test-id", Username: "test", Role: "admin"})
	h.HandleCreateToken(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	var resp ErrorResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected VALIDATION_ERROR, got %s", resp.Error.Code)
	}
}
