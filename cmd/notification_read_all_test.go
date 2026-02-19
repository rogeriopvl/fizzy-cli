package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/config"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestNotificationReadAllCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/test-account/notifications/bulk_reading" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusNoContent)
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

	cmd := notificationReadAllCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleReadAllNotifications(cmd); err != nil {
		t.Fatalf("handleReadAllNotifications failed: %v", err)
	}
}

func TestNotificationReadAllCommandAPIError(t *testing.T) {
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

	cmd := notificationReadAllCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleReadAllNotifications(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
}

func TestNotificationReadAllCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := notificationReadAllCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleReadAllNotifications(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
