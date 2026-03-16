package cmd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/config"
	"github.com/rogeriopvl/fizzy-cli/internal/testutil"
)

func TestWebhookCreateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/board-123/webhooks" {
			t.Errorf("expected /test-account/boards/board-123/webhooks, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]fizzy.CreateWebhookPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		webhookPayload := payload["webhook"]
		if webhookPayload.Name != "My Webhook" {
			t.Errorf("expected name 'My Webhook', got %s", webhookPayload.Name)
		}
		if webhookPayload.URL != "https://example.com/hook" {
			t.Errorf("expected URL 'https://example.com/hook', got %s", webhookPayload.URL)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fizzy.Webhook{
			ID:   "webhook-456",
			Name: "My Webhook",
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := webhookCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123", "--name", "My Webhook", "--url", "https://example.com/hook"})

	if err := handleCreateWebhook(cmd); err != nil {
		t.Fatalf("handleCreateWebhook failed: %v", err)
	}
}

func TestWebhookCreateCommandWithSelectedBoard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/selected-board/webhooks" {
			t.Errorf("expected /test-account/boards/selected-board/webhooks, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fizzy.Webhook{
			ID:   "webhook-456",
			Name: "My Webhook",
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "selected-board"},
	}

	cmd := webhookCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "", "--name", "My Webhook", "--url", "https://example.com/hook"})

	if err := handleCreateWebhook(cmd); err != nil {
		t.Fatalf("handleCreateWebhook failed: %v", err)
	}
}

func TestWebhookCreateCommandFlagOverridesSelectedBoard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/flag-board/webhooks" {
			t.Errorf("expected /test-account/boards/flag-board/webhooks, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fizzy.Webhook{
			ID:   "webhook-456",
			Name: "My Webhook",
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "selected-board"},
	}

	cmd := webhookCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "flag-board", "--name", "My Webhook", "--url", "https://example.com/hook"})

	if err := handleCreateWebhook(cmd); err != nil {
		t.Fatalf("handleCreateWebhook failed: %v", err)
	}
}

func TestWebhookCreateCommandWithActions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]fizzy.CreateWebhookPayload
		json.Unmarshal(body, &payload)

		webhookPayload := payload["webhook"]
		if len(webhookPayload.SubscribedActions) != 2 {
			t.Errorf("expected 2 actions, got %d", len(webhookPayload.SubscribedActions))
		}
		if webhookPayload.SubscribedActions[0] != "card.created" {
			t.Errorf("expected first action 'card.created', got %s", webhookPayload.SubscribedActions[0])
		}
		if webhookPayload.SubscribedActions[1] != "card.updated" {
			t.Errorf("expected second action 'card.updated', got %s", webhookPayload.SubscribedActions[1])
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fizzy.Webhook{
			ID:   "webhook-456",
			Name: "My Webhook",
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := webhookCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--board-id", "board-123",
		"--name", "My Webhook",
		"--url", "https://example.com/hook",
		"--actions", "card.created,card.updated",
	})

	if err := handleCreateWebhook(cmd); err != nil {
		t.Fatalf("handleCreateWebhook failed: %v", err)
	}
}

func TestWebhookCreateCommandNoBoard(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{},
	}

	cmd := webhookCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "", "--name", "My Webhook", "--url", "https://example.com/hook"})

	err := handleCreateWebhook(cmd)
	if err == nil {
		t.Errorf("expected error when no board specified")
	}
	if err.Error() != "no board specified: use --board-id or select a board with 'fizzy use'" {
		t.Errorf("expected no board error, got %v", err)
	}
}

func TestWebhookCreateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := webhookCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123", "--name", "My Webhook", "--url", "https://example.com/hook"})

	err := handleCreateWebhook(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "creating webhook: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestWebhookCreateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := webhookCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123", "--name", "My Webhook", "--url", "https://example.com/hook"})

	err := handleCreateWebhook(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
