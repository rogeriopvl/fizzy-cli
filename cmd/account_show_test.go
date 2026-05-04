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

func TestAccountShowCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/account/settings" {
			t.Errorf("expected /test-account/account/settings, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fizzy.Account{
			ID:                       "acc-123",
			Name:                     "37signals",
			CardsCount:               5,
			AutoPostponePeriodInDays: 30,
			CreatedAt:                "2025-12-05T19:36:35Z",
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := accountShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleShowAccount(cmd); err != nil {
		t.Fatalf("handleShowAccount failed: %v", err)
	}
}

func TestAccountShowCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := accountShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowAccount(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching account: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestAccountShowCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := accountShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowAccount(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
