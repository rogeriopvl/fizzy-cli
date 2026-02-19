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
)

func TestCardReactionListCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/123/reactions" {
			t.Errorf("expected /test-account/cards/123/reactions, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		reactions := []fizzy.Reaction{
			{
				ID:      "reaction-1",
				Content: "👍",
				Reacter: fizzy.User{ID: "user-1", Name: "Alice"},
			},
			{
				ID:      "reaction-2",
				Content: "🎉",
				Reacter: fizzy.User{ID: "user-2", Name: "Bob"},
			},
		}
		json.NewEncoder(w).Encode(reactions)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardReactionListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListCardReactions(cmd, "123"); err != nil {
		t.Fatalf("handleListCardReactions failed: %v", err)
	}
}

func TestCardReactionListCommandNoReactions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]fizzy.Reaction{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardReactionListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListCardReactions(cmd, "123"); err != nil {
		t.Fatalf("handleListCardReactions failed: %v", err)
	}
}

func TestCardReactionListCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := cardReactionListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListCardReactions(cmd, "not-a-number")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}

func TestCardReactionListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := cardReactionListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListCardReactions(cmd, "123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestCardReactionListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Card not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardReactionListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListCardReactions(cmd, "123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching reactions: unexpected status code 404: Card not found" {
		t.Errorf("expected API error, got %v", err)
	}
}
