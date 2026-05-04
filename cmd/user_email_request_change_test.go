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
)

func TestUserEmailRequestChangeCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/users/u-1/email_addresses" {
			t.Errorf("expected /test-account/users/u-1/email_addresses, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload fizzy.RequestEmailChangePayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}
		if payload.EmailAddress != "new@example.com" {
			t.Errorf("expected email_address=new@example.com, got %s", payload.EmailAddress)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userEmailRequestChangeCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--email", "new@example.com"})

	if err := handleRequestUserEmailChange(cmd, "u-1"); err != nil {
		t.Fatalf("handleRequestUserEmailChange failed: %v", err)
	}
}

func TestUserEmailRequestChangeCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userEmailRequestChangeCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--email", "new@example.com"})

	err := handleRequestUserEmailChange(cmd, "u-1")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "requesting email change: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestUserEmailRequestChangeCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := userEmailRequestChangeCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--email", "new@example.com"})

	err := handleRequestUserEmailChange(cmd, "u-1")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
