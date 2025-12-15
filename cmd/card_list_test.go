package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/config"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestCardListCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cards" {
			t.Errorf("expected /cards, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		boardIDs := r.URL.Query()["board_ids[]"]
		if len(boardIDs) == 0 || boardIDs[0] != "board-123" {
			t.Errorf("expected board_ids[]=board-123 in query, got %v", boardIDs)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := []api.Card{
			{
				ID:        "card-123",
				Number:    1,
				Title:     "Implement feature",
				Status:    "in_progress",
				CreatedAt: "2025-01-01T00:00:00Z",
			},
			{
				ID:        "card-456",
				Number:    2,
				Title:     "Fix bug",
				Status:    "todo",
				CreatedAt: "2025-01-02T00:00:00Z",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := cardListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListCards(cmd); err != nil {
		t.Fatalf("handleListCards failed: %v", err)
	}
}

func TestCardListCommandNoCards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]api.Card{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := cardListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListCards(cmd); err != nil {
		t.Fatalf("handleListCards failed: %v", err)
	}
}

func TestCardListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := cardListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListCards(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching cards: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestCardListCommandNoBoard(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: ""},
	}

	cmd := cardListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListCards(cmd)
	if err == nil {
		t.Errorf("expected error when board not selected")
	}
	if err.Error() != "no board selected" {
		t.Errorf("expected 'no board selected' error, got %v", err)
	}
}

func TestCardListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := cardListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListCards(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
