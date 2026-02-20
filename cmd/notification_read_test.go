package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/config"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestNotificationReadCommand(t *testing.T) {
	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if r.URL.Path == "/test-account/notifications/notif-123/reading" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.URL.Path == "/test-account/notifications/notif-123" && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			response := fizzy.Notification{
				ID:        "notif-123",
				Read:      true,
				ReadAt:    "2025-01-01T00:00:00Z",
				CreatedAt: "2025-01-01T00:00:00Z",
				Title:     "Test Notification",
				Body:      "This is a test notification",
				Creator: fizzy.User{
					ID:        "user-123",
					Name:      "David Heinemeier Hansson",
					Email:     "david@example.com",
					Role:      "owner",
					Active:    true,
					CreatedAt: "2025-12-05T19:36:35.401Z",
				},
				Card: fizzy.CardReference{
					ID:     "card-123",
					Title:  "Test Card",
					Status: "published",
					URL:    "http://fizzy.localhost:3006/897362094/cards/1",
				},
				URL: "http://fizzy.localhost:3006/897362094/notifications/notif-123",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := notificationReadCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleReadNotification(cmd, "notif-123"); err != nil {
		t.Fatalf("handleReadNotification failed: %v", err)
	}

	if requestCount != 2 {
		t.Errorf("expected 2 requests (POST and GET), got %d", requestCount)
	}
}

func TestNotificationReadCommandNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Notification not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "board-123", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{SelectedBoard: "board-123"},
	}

	cmd := notificationReadCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleReadNotification(cmd, "notif-invalid")
	if err == nil {
		t.Errorf("expected error for invalid notification")
	}
	if err.Error() != "notification not found" {
		t.Errorf("expected 'notification not found' error, got %v", err)
	}
}

func TestNotificationReadCommandAPIError(t *testing.T) {
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

	cmd := notificationReadCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleReadNotification(cmd, "notif-123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
}

func TestNotificationReadCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := notificationReadCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleReadNotification(cmd, "notif-123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
