package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
	"github.com/spf13/cobra"
)

func TestCardUpdateCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/1" {
			t.Errorf("expected /test-account/cards/1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := fizzy.Card{
			ID:           "card-123",
			Number:       1,
			Title:        "Updated card title",
			Status:       "published",
			Description:  "Updated description",
			Tags:         []string{"updated-tag"},
			Golden:       false,
			CreatedAt:    "2025-01-01T00:00:00Z",
			LastActiveAt: "2025-01-15T10:30:00Z",
			URL:          "https://example.com/card/1",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--title", "Updated card title",
		"--description", "Updated description",
		"--status", "published",
		"--tag-id", "updated-tag",
	})

	if err := handleUpdateCard(cmd, "1"); err != nil {
		t.Fatalf("handleUpdateCard failed: %v", err)
	}
}

func TestCardUpdateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Card not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--title", "Updated card",
	})

	err := handleUpdateCard(cmd, "999")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "updating card: unexpected status code 404: Card not found" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestCardUpdateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := cardUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--title", "Updated card",
	})

	err := handleUpdateCard(cmd, "1")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestCardUpdateCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := cardUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--title", "Updated card",
	})

	err := handleUpdateCard(cmd, "not-a-number")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}

func TestCardUpdateCommandNoFlags(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{Client: client}

	// Create a fresh command to avoid flag pollution from other tests
	cmd := &cobra.Command{
		Use:  "update <card_number>",
		Args: cobra.ExactArgs(1),
	}
	cmd.Flags().String("title", "", "Card title")
	cmd.Flags().String("description", "", "Card description")
	cmd.Flags().String("status", "", "Card status")
	cmd.Flags().StringSlice("tag-id", []string{}, "Tag ID")
	cmd.Flags().String("last-active-at", "", "Last active timestamp")

	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleUpdateCard(cmd, "1")
	if err == nil {
		t.Errorf("expected error when no flags are provided")
	}
	if err.Error() != "must provide at least one flag to update (--title, --description, --status, --tag-id, or --last-active-at)" {
		t.Errorf("expected 'no flags' error, got %v", err)
	}
}
