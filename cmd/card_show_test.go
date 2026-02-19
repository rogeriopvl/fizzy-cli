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

func TestCardShowCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/1" {
			t.Errorf("expected /test-account/cards/1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := fizzy.Card{
			ID:           "card-123",
			Number:       1,
			Title:        "Implement feature",
			Status:       "in_progress",
			Description:  "This is a test card",
			Tags:         []string{"feature", "backend"},
			Golden:       false,
		Closed:       false,
		Column:       &fizzy.Column{ID: "col-1", Name: "In Progress"},
		Steps:        []fizzy.Step{{ID: "step-1", Content: "First step", Completed: true}},
		CreatedAt:    "2025-01-01T00:00:00Z",
			LastActiveAt: "2025-01-15T10:30:00Z",
			URL:          "https://example.com/card/1",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleShowCard(cmd, "1"); err != nil {
		t.Fatalf("handleShowCard failed: %v", err)
	}
}

func TestCardShowCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Card not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowCard(cmd, "999")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching card: unexpected status code 404: Card not found" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestCardShowCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := cardShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowCard(cmd, "1")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
