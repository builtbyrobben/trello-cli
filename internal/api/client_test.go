package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient_Defaults(t *testing.T) {
	t.Parallel()

	c := NewClient("test-key")
	if c.baseURL != "https://api.trello.com/1" {
		t.Errorf("unexpected baseURL: %s", c.baseURL)
	}

	if c.userAgent != "trello-cli/1.0" {
		t.Errorf("unexpected userAgent: %s", c.userAgent)
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	t.Parallel()

	c := NewClient("test-key",
		WithBaseURL("https://custom.api.com"),
		WithUserAgent("custom/2.0"),
	)

	if c.baseURL != "https://custom.api.com" {
		t.Errorf("unexpected baseURL: %s", c.baseURL)
	}

	if c.userAgent != "custom/2.0" {
		t.Errorf("unexpected userAgent: %s", c.userAgent)
	}
}

func TestQueryAuth(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		token := r.URL.Query().Get("token")

		if key != "mykey" {
			t.Errorf("expected query param key=mykey, got %s", key)
		}

		if token != "mytoken" {
			t.Errorf("expected query param token=mytoken, got %s", token)
		}

		// Ensure no Authorization header
		if auth := r.Header.Get("Authorization"); auth != "" {
			t.Errorf("expected no Authorization header, got %s", auth)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer srv.Close()

	c := NewClient("",
		WithBaseURL(srv.URL),
		WithQueryAuth(map[string]string{
			"key":   "mykey",
			"token": "mytoken",
		}),
	)

	var result map[string]string

	err := c.Get(context.Background(), "/test", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["ok"] != "true" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestGet_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "123"})
	}))
	defer srv.Close()

	c := NewClient("test", WithBaseURL(srv.URL))

	var result map[string]string

	err := c.Get(context.Background(), "/test", &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["id"] != "123" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestGet_APIError(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found"})
	}))
	defer srv.Close()

	c := NewClient("test", WithBaseURL(srv.URL))

	var result map[string]string

	err := c.Get(context.Background(), "/missing", &result)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", apiErr.StatusCode)
	}
}

func TestPost_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)

		if body["name"] != "test card" {
			t.Errorf("unexpected body name: %s", body["name"])
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "456", "name": "test card"})
	}))
	defer srv.Close()

	c := NewClient("test", WithBaseURL(srv.URL))

	var result map[string]string

	err := c.Post(context.Background(), "/cards", map[string]string{"name": "test card"}, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["id"] != "456" {
		t.Errorf("unexpected id: %s", result["id"])
	}
}

func TestDelete_Success(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := NewClient("test", WithBaseURL(srv.URL))

	err := c.Delete(context.Background(), "/cards/123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
