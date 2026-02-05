package cmd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestCardReactionCreateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cards/123/reactions" {
			t.Errorf("expected /cards/123/reactions, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]map[string]string
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		reactionPayload := payload["reaction"]
		if reactionPayload["content"] != "👍" {
			t.Errorf("expected content '👍', got %s", reactionPayload["content"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := api.Reaction{
			ID:      "reaction-123",
			Content: "👍",
			Reacter: api.User{ID: "user-1", Name: "John Doe"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardReactionCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleCreateCardReaction(cmd, "123", "👍"); err != nil {
		t.Fatalf("handleCreateCardReaction failed: %v", err)
	}
}

func TestCardReactionCreateCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := cardReactionCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleCreateCardReaction(cmd, "not-a-number", "👍")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}

func TestCardReactionCreateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := cardReactionCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleCreateCardReaction(cmd, "123", "👍")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestCardReactionCreateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Card not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardReactionCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleCreateCardReaction(cmd, "123", "👍")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "creating reaction: unexpected status code 404: Card not found" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestCardReactionCreateCommandDifferentEmoji(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]map[string]string
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		reactionPayload := payload["reaction"]
		if reactionPayload["content"] != "🎉" {
			t.Errorf("expected content '🎉', got %s", reactionPayload["content"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := api.Reaction{
			ID:      "reaction-456",
			Content: "🎉",
			Reacter: api.User{ID: "user-2", Name: "Jane Doe"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardReactionCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleCreateCardReaction(cmd, "123", "🎉"); err != nil {
		t.Fatalf("handleCreateCardReaction failed: %v", err)
	}
}
