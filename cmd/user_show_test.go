package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestUserShowCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/user-123" {
			t.Errorf("expected /users/user-123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := api.User{
			ID:        "user-123",
			Name:      "John Doe",
			Email:     "john@example.com",
			Role:      "admin",
			Active:    true,
			CreatedAt: "2025-01-01T00:00:00Z",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleShowUser(cmd, "user-123"); err != nil {
		t.Fatalf("handleShowUser failed: %v", err)
	}
}

func TestUserShowCommandNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowUser(cmd, "nonexistent-user")
	if err == nil {
		t.Errorf("expected error for user not found")
	}
	if err.Error() != "fetching user: unexpected status code 404: User not found" {
		t.Errorf("expected user not found error, got %v", err)
	}
}

func TestUserShowCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowUser(cmd, "user-123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching user: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestUserShowCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := userShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowUser(cmd, "user-123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
