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

func TestColumnShowCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/boards/board-123/columns/col-456" {
			t.Errorf("expected /boards/board-123/columns/col-456, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := api.Column{
			ID:        "col-456",
			Name:      "In Progress",
			CreatedAt: "2025-01-01T00:00:00Z",
			Color: api.ColorObject{
				Name:  "Lime",
				Value: api.Lime,
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleShowColumnDetails(cmd, "col-456"); err != nil {
		t.Fatalf("handleShowColumnDetails failed: %v", err)
	}
}

func TestColumnShowCommandNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Column not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowColumnDetails(cmd, "nonexistent-col")
	if err == nil {
		t.Errorf("expected error for column not found")
	}
	if err.Error() != "fetching column: unexpected status code 404: Column not found" {
		t.Errorf("expected column not found error, got %v", err)
	}
}

func TestColumnShowCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowColumnDetails(cmd, "col-456")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching column: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestColumnShowCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := columnShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowColumnDetails(cmd, "col-456")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
