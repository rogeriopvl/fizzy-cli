package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/testutil"
)

func TestUserEmailConfirmChangeCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/users/u-1/email_addresses/tok-xyz/confirmation" {
			t.Errorf("expected /test-account/users/u-1/email_addresses/tok-xyz/confirmation, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userEmailConfirmChangeCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--token", "tok-xyz"})

	if err := handleConfirmUserEmailChange(cmd, "u-1"); err != nil {
		t.Fatalf("handleConfirmUserEmailChange failed: %v", err)
	}
}

func TestUserEmailConfirmChangeCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userEmailConfirmChangeCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--token", "tok-xyz"})

	err := handleConfirmUserEmailChange(cmd, "u-1")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "confirming email change: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestUserEmailConfirmChangeCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := userEmailConfirmChangeCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--token", "tok-xyz"})

	err := handleConfirmUserEmailChange(cmd, "u-1")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
