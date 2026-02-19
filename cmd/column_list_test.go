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

func TestColumnListCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/board-123/columns" {
			t.Errorf("expected /test-account/boards/board-123/columns, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("expected Accept: application/json, got %s", r.Header.Get("Accept"))
		}

		w.Header().Set("Content-Type", "application/json")
		response := []fizzy.Column{
			{
				ID:   "col-123",
				Name: "Todo",
				Color: fizzy.ColorObject{
					Name:  "Blue",
					Value: "var(--color-card-default)",
				},
				CreatedAt: "2025-01-01T00:00:00Z",
			},
			{
				ID:   "col-456",
				Name: "In Progress",
				Color: fizzy.ColorObject{
					Name:  "Lime",
					Value: "var(--color-card-4)",
				},
				CreatedAt: "2025-01-02T00:00:00Z",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListColumns(cmd); err != nil {
		t.Fatalf("handleListColumns failed: %v", err)
	}
}

func TestColumnListCommandNoColumns(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Column{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListColumns(cmd); err != nil {
		t.Fatalf("handleListColumns failed: %v", err)
	}
}

func TestColumnListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListColumns(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching columns: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestColumnListCommandNoBoard(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := columnListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListColumns(cmd)
	if err == nil {
		t.Errorf("expected error when board not selected")
	}
	if err.Error() != "fetching columns: no board selected: use SetBoard or WithBoard when creating the client" {
		t.Errorf("expected board error, got %v", err)
	}
}

func TestColumnListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := columnListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListColumns(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
