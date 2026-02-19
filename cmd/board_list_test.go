package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestBoardListCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards" {
			t.Errorf("expected /test-account/boards, got %s", r.URL.Path)
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
		response := []fizzy.Board{
			{
				ID:        "board-123",
				Name:      "Project Alpha",
				AllAccess: true,
				CreatedAt: "2025-01-01T00:00:00Z",
			},
			{
				ID:        "board-456",
				Name:      "Project Beta",
				AllAccess: false,
				CreatedAt: "2025-01-02T00:00:00Z",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListBoards(cmd); err != nil {
		t.Fatalf("handleListBoards failed: %v", err)
	}
}

func TestBoardListCommandNoBoards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Board{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListBoards(cmd); err != nil {
		t.Fatalf("handleListBoards failed: %v", err)
	}
}

func TestBoardListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListBoards(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching boards: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestBoardListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := boardListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListBoards(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
