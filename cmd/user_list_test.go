package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestUserListCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/users" {
			t.Errorf("expected /test-account/users, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("missing Authorization header")
		}
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := []fizzy.User{
			{
				ID:        "user-123",
				Name:      "John Doe",
				Email:     "john@example.com",
				Role:      "admin",
				Active:    true,
				CreatedAt: "2025-01-01T00:00:00Z",
			},
			{
				ID:        "user-456",
				Name:      "Jane Smith",
				Email:     "jane@example.com",
				Role:      "member",
				Active:    false,
				CreatedAt: "2025-01-02T00:00:00Z",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.SetOut(&bytes.Buffer{})

	if err := handleListUsers(cmd); err != nil {
		t.Fatalf("handleListUsers failed: %v", err)
	}
}

func TestUserListCommandNoUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.User{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.SetOut(&bytes.Buffer{})

	if err := handleListUsers(cmd); err != nil {
		t.Fatalf("handleListUsers failed: %v", err)
	}
}

func TestUserListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	cmd.SetOut(&bytes.Buffer{})

	err := handleListUsers(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching users: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestUserListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := userListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	cmd.SetOut(&bytes.Buffer{})

	err := handleListUsers(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
