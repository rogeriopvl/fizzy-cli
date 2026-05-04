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

func TestBoardEntropyCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/boards/board-123/entropy" {
			t.Errorf("expected /test-account/boards/board-123/entropy, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]fizzy.EntropyPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}
		if payload["board"].AutoPostponePeriodInDays != 60 {
			t.Errorf("expected auto_postpone_period_in_days 60, got %d", payload["board"].AutoPostponePeriodInDays)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fizzy.Board{
			ID:                       "board-123",
			Name:                     "Mobile",
			AutoPostponePeriodInDays: 60,
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardEntropyCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--auto-postpone-days", "60"})

	if err := handleBoardEntropy(cmd, "board-123"); err != nil {
		t.Fatalf("handleBoardEntropy failed: %v", err)
	}
}

func TestBoardEntropyCommandMissingFlag(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := &cobra.Command{Use: "entropy <board_id>", Args: cobra.ExactArgs(1)}
	cmd.Flags().Int("auto-postpone-days", 0, "Auto-postpone period in days (required)")
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleBoardEntropy(cmd, "board-123")
	if err == nil {
		t.Errorf("expected error when --auto-postpone-days not provided")
	}
	if err.Error() != "--auto-postpone-days is required" {
		t.Errorf("expected required-flag error, got %v", err)
	}
}

func TestBoardEntropyCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := boardEntropyCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--auto-postpone-days", "60"})

	err := handleBoardEntropy(cmd, "board-123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "updating board entropy: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestBoardEntropyCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := boardEntropyCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--auto-postpone-days", "60"})

	err := handleBoardEntropy(cmd, "board-123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
