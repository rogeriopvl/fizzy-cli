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
	"github.com/spf13/cobra"
)

func TestWebhookUpdateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/board-123/webhooks/webhook-456" {
			t.Errorf("expected /test-account/boards/board-123/webhooks/webhook-456, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]fizzy.UpdateWebhookPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		webhookPayload := payload["webhook"]
		if webhookPayload.Name != "Updated Hook" {
			t.Errorf("expected name 'Updated Hook', got %s", webhookPayload.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fizzy.Webhook{
			ID:   "webhook-456",
			Name: "Updated Hook",
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := webhookUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123", "--name", "Updated Hook"})

	if err := handleUpdateWebhook(cmd, "webhook-456"); err != nil {
		t.Fatalf("handleUpdateWebhook failed: %v", err)
	}
}

func TestWebhookUpdateCommandWithActions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]fizzy.UpdateWebhookPayload
		json.Unmarshal(body, &payload)

		webhookPayload := payload["webhook"]
		if len(webhookPayload.SubscribedActions) != 2 {
			t.Errorf("expected 2 actions, got %d", len(webhookPayload.SubscribedActions))
		}
		if webhookPayload.SubscribedActions[0] != "card.created" {
			t.Errorf("expected first action 'card.created', got %s", webhookPayload.SubscribedActions[0])
		}
		if webhookPayload.SubscribedActions[1] != "card.closed" {
			t.Errorf("expected second action 'card.closed', got %s", webhookPayload.SubscribedActions[1])
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fizzy.Webhook{
			ID:   "webhook-456",
			Name: "Hook",
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := webhookUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123", "--actions", "card.created,card.closed"})

	if err := handleUpdateWebhook(cmd, "webhook-456"); err != nil {
		t.Fatalf("handleUpdateWebhook failed: %v", err)
	}
}

func TestWebhookUpdateCommandWithSelectedBoard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/selected-board/webhooks/webhook-456" {
			t.Errorf("expected /test-account/boards/selected-board/webhooks/webhook-456, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fizzy.Webhook{
			ID:   "webhook-456",
			Name: "Updated Hook",
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "selected-board"},
	}

	cmd := webhookUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "", "--name", "Updated Hook"})

	if err := handleUpdateWebhook(cmd, "webhook-456"); err != nil {
		t.Fatalf("handleUpdateWebhook failed: %v", err)
	}
}

func TestWebhookUpdateCommandNoFlags(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := &cobra.Command{
		Use:  "update <webhook_id>",
		Args: cobra.ExactArgs(1),
	}
	cmd.Flags().String("board-id", "", "Board ID")
	cmd.Flags().String("name", "", "Webhook name")
	cmd.Flags().StringSlice("actions", nil, "Subscribed actions")

	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleUpdateWebhook(cmd, "webhook-456")
	if err == nil {
		t.Errorf("expected error when no flags provided")
	}
	if err.Error() != "at least one flag must be provided (--name or --actions)" {
		t.Errorf("expected flag requirement error, got %v", err)
	}
}

func TestWebhookUpdateCommandNoBoard(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{},
	}

	cmd := webhookUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "", "--name", "Updated Hook"})

	err := handleUpdateWebhook(cmd, "webhook-456")
	if err == nil {
		t.Errorf("expected error when no board specified")
	}
	if err.Error() != "no board specified: use --board-id or select a board with 'fizzy use'" {
		t.Errorf("expected no board error, got %v", err)
	}
}

func TestWebhookUpdateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := webhookUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123", "--name", "Updated Hook"})

	err := handleUpdateWebhook(cmd, "webhook-456")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "updating webhook: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestWebhookUpdateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := webhookUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123", "--name", "Updated Hook"})

	err := handleUpdateWebhook(cmd, "webhook-456")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
