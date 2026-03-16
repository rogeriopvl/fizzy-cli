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

func TestWebhookActivateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/board-123/webhooks/webhook-456/activation" {
			t.Errorf("expected /test-account/boards/board-123/webhooks/webhook-456/activation, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fizzy.Webhook{
			ID:     "webhook-456",
			Name:   "Deploy Hook",
			Active: true,
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := webhookActivateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123"})

	if err := handleActivateWebhook(cmd, "webhook-456"); err != nil {
		t.Fatalf("handleActivateWebhook failed: %v", err)
	}
}

func TestWebhookActivateCommandWithSelectedBoard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/selected-board/webhooks/webhook-456/activation" {
			t.Errorf("expected /test-account/boards/selected-board/webhooks/webhook-456/activation, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fizzy.Webhook{
			ID:     "webhook-456",
			Name:   "Deploy Hook",
			Active: true,
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "selected-board"},
	}

	cmd := webhookActivateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", ""})

	if err := handleActivateWebhook(cmd, "webhook-456"); err != nil {
		t.Fatalf("handleActivateWebhook failed: %v", err)
	}
}

func TestWebhookActivateCommandNoBoard(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{},
	}

	cmd := webhookActivateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", ""})

	err := handleActivateWebhook(cmd, "webhook-456")
	if err == nil {
		t.Errorf("expected error when no board specified")
	}
	if err.Error() != "no board specified: use --board-id or select a board with 'fizzy use'" {
		t.Errorf("expected no board error, got %v", err)
	}
}

func TestWebhookActivateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := webhookActivateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123"})

	err := handleActivateWebhook(cmd, "webhook-456")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "activating webhook: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestWebhookActivateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := webhookActivateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123"})

	err := handleActivateWebhook(cmd, "webhook-456")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
