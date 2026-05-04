package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/testutil"
)

func TestCardPinCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/123/pin" {
			t.Errorf("expected /test-account/cards/123/pin, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardPinCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handlePinCard(cmd, "123"); err != nil {
		t.Fatalf("handlePinCard failed: %v", err)
	}
}

func TestCardPinCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := cardPinCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handlePinCard(cmd, "not-a-number")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
}

func TestCardPinCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := cardPinCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handlePinCard(cmd, "123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestCardPinCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardPinCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handlePinCard(cmd, "123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "pinning card: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}
