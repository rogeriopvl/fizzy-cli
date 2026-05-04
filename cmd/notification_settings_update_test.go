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
	"github.com/rogeriopvl/fizzy-cli/internal/testutil"
	"github.com/spf13/cobra"
)

func TestNotificationSettingsUpdateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/notifications/settings" {
			t.Errorf("expected /test-account/notifications/settings, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]fizzy.UpdateNotificationSettingsPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}
		if payload["user_settings"].BundleEmailFrequency != "daily" {
			t.Errorf("expected bundle_email_frequency=daily, got %s", payload["user_settings"].BundleEmailFrequency)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := notificationSettingsUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--bundle-email-frequency", "daily"})

	if err := handleUpdateNotificationSettings(cmd); err != nil {
		t.Fatalf("handleUpdateNotificationSettings failed: %v", err)
	}
}

func TestNotificationSettingsUpdateCommandMissingFlag(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := &cobra.Command{}
	cmd.Flags().String("bundle-email-frequency", "", "")
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleUpdateNotificationSettings(cmd)
	if err == nil {
		t.Errorf("expected error when no flags provided")
	}
}

func TestNotificationSettingsUpdateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := notificationSettingsUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--bundle-email-frequency", "daily"})

	err := handleUpdateNotificationSettings(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "updating notification settings: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestNotificationSettingsUpdateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := notificationSettingsUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--bundle-email-frequency", "daily"})

	err := handleUpdateNotificationSettings(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
