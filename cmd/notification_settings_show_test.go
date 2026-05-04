package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/testutil"
)

func TestNotificationSettingsShowCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/notifications/settings" {
			t.Errorf("expected /test-account/notifications/settings, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fizzy.NotificationSettings{BundleEmailFrequency: "daily"})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := notificationSettingsShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleShowNotificationSettings(cmd); err != nil {
		t.Fatalf("handleShowNotificationSettings failed: %v", err)
	}
}

func TestNotificationSettingsShowCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := notificationSettingsShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowNotificationSettings(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching notification settings: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestNotificationSettingsShowCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := notificationSettingsShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowNotificationSettings(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
