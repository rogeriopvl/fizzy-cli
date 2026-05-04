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

func TestAccountEntropyCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/account/entropy" {
			t.Errorf("expected /test-account/account/entropy, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]fizzy.EntropyPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}
		if payload["entropy"].AutoPostponePeriodInDays != 30 {
			t.Errorf("expected auto_postpone_period_in_days 30, got %d", payload["entropy"].AutoPostponePeriodInDays)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fizzy.Account{
			ID:                       "acc-123",
			Name:                     "37signals",
			AutoPostponePeriodInDays: 30,
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := accountEntropyCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--auto-postpone-days", "30"})

	if err := handleAccountEntropy(cmd); err != nil {
		t.Fatalf("handleAccountEntropy failed: %v", err)
	}
}

func TestAccountEntropyCommandMissingFlag(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := &cobra.Command{Use: "entropy"}
	cmd.Flags().Int("auto-postpone-days", 0, "Auto-postpone period in days (required)")
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleAccountEntropy(cmd)
	if err == nil {
		t.Errorf("expected error when --auto-postpone-days not provided")
	}
	if err.Error() != "--auto-postpone-days is required" {
		t.Errorf("expected required-flag error, got %v", err)
	}
}

func TestAccountEntropyCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := accountEntropyCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--auto-postpone-days", "30"})

	err := handleAccountEntropy(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "updating account entropy: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestAccountEntropyCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := accountEntropyCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--auto-postpone-days", "30"})

	err := handleAccountEntropy(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
