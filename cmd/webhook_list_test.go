package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/config"
	"github.com/rogeriopvl/fizzy-cli/internal/testutil"
)

func TestWebhookListCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/board-123/webhooks" {
			t.Errorf("expected /test-account/boards/board-123/webhooks, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := []fizzy.Webhook{
			{
				ID:     "webhook-123",
				Name:   "Deploy Hook",
				Active: true,
			},
			{
				ID:     "webhook-456",
				Name:   "CI Hook",
				Active: false,
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := webhookListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123"})

	if err := handleListWebhooks(cmd); err != nil {
		t.Fatalf("handleListWebhooks failed: %v", err)
	}
}

func TestWebhookListCommandWithSelectedBoard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/selected-board/webhooks" {
			t.Errorf("expected /test-account/boards/selected-board/webhooks, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Webhook{
			{ID: "webhook-123", Name: "Hook", Active: true},
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "selected-board"},
	}

	cmd := webhookListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", ""})

	if err := handleListWebhooks(cmd); err != nil {
		t.Fatalf("handleListWebhooks failed: %v", err)
	}
}

func TestWebhookListCommandFlagOverridesSelectedBoard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/flag-board/webhooks" {
			t.Errorf("expected /test-account/boards/flag-board/webhooks, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Webhook{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "selected-board"},
	}

	cmd := webhookListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "flag-board"})

	if err := handleListWebhooks(cmd); err != nil {
		t.Fatalf("handleListWebhooks failed: %v", err)
	}
}

func TestWebhookListCommandEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.Webhook{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := webhookListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123"})

	if err := handleListWebhooks(cmd); err != nil {
		t.Fatalf("handleListWebhooks failed: %v", err)
	}
}

func TestWebhookListCommandNoBoard(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{},
	}

	cmd := webhookListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", ""})

	err := handleListWebhooks(cmd)
	if err == nil {
		t.Errorf("expected error when no board specified")
	}
	if err.Error() != "no board specified: use --board-id or select a board with 'fizzy use'" {
		t.Errorf("expected no board error, got %v", err)
	}
}

func TestWebhookListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := webhookListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123"})

	err := handleListWebhooks(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching webhooks: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestWebhookListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := webhookListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123"})

	err := handleListWebhooks(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
