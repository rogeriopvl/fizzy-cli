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

func TestWebhookDeliveryListCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/board-123/webhooks/wh-1/deliveries" {
			t.Errorf("expected /test-account/boards/board-123/webhooks/wh-1/deliveries, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.WebhookDelivery{
			{
				ID:        "del-1",
				State:     "delivered",
				CreatedAt: "2026-03-25T15:11:04Z",
				Response:  &fizzy.WebhookDeliveryResponse{Code: 200},
				Event:     fizzy.WebhookDeliveryEvent{Action: "card.created"},
			},
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := webhookDeliveryListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123"})

	if err := handleListWebhookDeliveries(cmd, "wh-1"); err != nil {
		t.Fatalf("handleListWebhookDeliveries failed: %v", err)
	}
}

func TestWebhookDeliveryListCommandWithSelectedBoard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/selected-board/webhooks/wh-1/deliveries" {
			t.Errorf("expected /test-account/boards/selected-board/webhooks/wh-1/deliveries, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]fizzy.WebhookDelivery{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "selected-board"},
	}

	cmd := webhookDeliveryListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", ""})

	if err := handleListWebhookDeliveries(cmd, "wh-1"); err != nil {
		t.Fatalf("handleListWebhookDeliveries failed: %v", err)
	}
}

func TestWebhookDeliveryListCommandNoBoard(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{Client: client, Config: &config.Config{}}

	cmd := webhookDeliveryListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", ""})

	err := handleListWebhookDeliveries(cmd, "wh-1")
	if err == nil {
		t.Errorf("expected error when no board specified")
	}
	if err.Error() != "no board specified: use --board-id or select a board with 'fizzy use'" {
		t.Errorf("expected no-board error, got %v", err)
	}
}

func TestWebhookDeliveryListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := webhookDeliveryListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123"})

	err := handleListWebhookDeliveries(cmd, "wh-1")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching webhook deliveries: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestWebhookDeliveryListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := webhookDeliveryListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board-id", "board-123"})

	err := handleListWebhookDeliveries(cmd, "wh-1")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
