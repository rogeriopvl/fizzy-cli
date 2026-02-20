package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestNotificationUnreadCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/notifications/notif-123/reading" {
			t.Errorf("expected /test-account/notifications/notif-123/reading, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := notificationUnreadCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleUnreadNotification(cmd, "notif-123"); err != nil {
		t.Fatalf("handleUnreadNotification failed: %v", err)
	}
}

func TestNotificationUnreadCommandNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Notification not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := notificationUnreadCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleUnreadNotification(cmd, "notif-invalid")
	if err == nil {
		t.Errorf("expected error for invalid notification")
	}
	if err.Error() != "notification not found" {
		t.Errorf("expected 'notification not found' error, got %v", err)
	}
}

func TestNotificationUnreadCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := notificationUnreadCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleUnreadNotification(cmd, "notif-123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "marking notification as unread: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestNotificationUnreadCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := notificationUnreadCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleUnreadNotification(cmd, "notif-123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
