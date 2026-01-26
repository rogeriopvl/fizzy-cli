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

func TestBoardShowCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/boards/board-123" {
			t.Errorf("expected /boards/board-123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := api.Board{
			ID:        "board-123",
			Name:      "Project Alpha",
			AllAccess: true,
			CreatedAt: "2025-01-01T00:00:00Z",
			Creator: api.User{
				ID:   "user-123",
				Name: "John Doe",
				Role: "owner",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleShowBoardDetails(cmd, "board-123"); err != nil {
		t.Fatalf("handleShowBoardDetails failed: %v", err)
	}
}

func TestBoardShowCommandNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Board not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowBoardDetails(cmd, "nonexistent-board")
	if err == nil {
		t.Errorf("expected error for board not found")
	}
	if err.Error() != "fetching board: unexpected status code 404: Board not found" {
		t.Errorf("expected board not found error, got %v", err)
	}
}

func TestBoardShowCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowBoardDetails(cmd, "board-123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching board: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestBoardShowCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := boardShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowBoardDetails(cmd, "board-123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
