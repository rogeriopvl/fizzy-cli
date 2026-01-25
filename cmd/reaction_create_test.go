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

func TestReactionCreateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cards/123/comments/comment-456/reactions" {
			t.Errorf("expected /cards/123/comments/comment-456/reactions, got %s", r.URL.Path)
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
			ID:      "reaction-789",
			Content: "👍",
			Reacter: api.User{ID: "user-1", Name: "John Doe"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := reactionCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleCreateReaction(cmd, "123", "comment-456", "👍"); err != nil {
		t.Fatalf("handleCreateReaction failed: %v", err)
	}
}

func TestReactionCreateCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := reactionCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleCreateReaction(cmd, "not-a-number", "comment-456", "👍")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}

func TestReactionCreateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := reactionCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleCreateReaction(cmd, "123", "comment-456", "👍")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestReactionCreateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Comment not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := reactionCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleCreateReaction(cmd, "123", "comment-456", "👍")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "creating reaction: unexpected status code 404: Comment not found" {
		t.Errorf("expected API error, got %v", err)
	}
}
