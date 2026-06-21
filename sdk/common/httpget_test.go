package common

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetJSON_SendsGETWithAcceptAndNoBearerByDefault(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Errorf("Accept = %q, want application/json", got)
		}
		if got := r.Header.Get("Authorization"); got != "" {
			t.Errorf("Authorization = %q, want empty (no bearer)", got)
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	body, err := GetJSON(context.Background(), srv.Client(), srv.URL, "", 4<<20)
	if err != nil {
		t.Fatalf("GetJSON: %v", err)
	}
	if string(body) != `{"ok":true}` {
		t.Errorf("body = %q, want {\"ok\":true}", body)
	}
}

func TestGetJSON_AttachesBearerWhenSet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer tok-123" {
			t.Errorf("Authorization = %q, want Bearer tok-123", got)
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	if _, err := GetJSON(context.Background(), srv.Client(), srv.URL, "tok-123", 4<<20); err != nil {
		t.Fatalf("GetJSON: %v", err)
	}
}

func TestGetJSON_TruncatesAtLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("a", 1000)))
	}))
	defer srv.Close()

	body, err := GetJSON(context.Background(), srv.Client(), srv.URL, "", 10)
	if err != nil {
		t.Fatalf("GetJSON: %v", err)
	}
	if len(body) != 10 {
		t.Errorf("len(body) = %d, want 10 (LimitReader cap)", len(body))
	}
}

func TestGetJSON_Non2xxIsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := GetJSON(context.Background(), srv.Client(), srv.URL, "", 4<<20)
	if err == nil {
		t.Fatal("GetJSON: want error on HTTP 500, got nil")
	}
	if !strings.Contains(err.Error(), "status 500") {
		t.Errorf("error = %v, want to contain 'status 500'", err)
	}
}

func TestNoRedirectClient_DoesNotFollowRedirects(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/start" {
			http.Redirect(w, r, "/elsewhere", http.StatusFound)
			return
		}
		t.Errorf("redirect was followed to %s (NoRedirectClient must stop)", r.URL.Path)
	}))
	defer srv.Close()

	resp, err := NoRedirectClient(5 * time.Second).Get(srv.URL + "/start")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusFound {
		t.Errorf("status = %d, want 302 (redirect not followed)", resp.StatusCode)
	}
}
