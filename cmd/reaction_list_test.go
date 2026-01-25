package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestReactionListCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cards/123/comments/comment-456/reactions" {
			t.Errorf("expected /cards/123/comments/comment-456/reactions, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := []api.Reaction{
			{
				ID:      "reaction-1",
				Content: "👍",
				Reacter: api.User{ID: "user-1", Name: "John Doe"},
			},
			{
				ID:      "reaction-2",
				Content: "❤️",
				Reacter: api.User{ID: "user-2", Name: "Jane Doe"},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := reactionListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListReactions(cmd, "123", "comment-456"); err != nil {
		t.Fatalf("handleListReactions failed: %v", err)
	}
}

func TestReactionListCommandNoReactions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]api.Reaction{})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := reactionListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListReactions(cmd, "123", "comment-456"); err != nil {
		t.Fatalf("handleListReactions failed: %v", err)
	}
}

func TestReactionListCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := reactionListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListReactions(cmd, "not-a-number", "comment-456")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}

func TestReactionListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := reactionListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListReactions(cmd, "123", "comment-456")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestReactionListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Comment not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := reactionListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListReactions(cmd, "123", "comment-456")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching reactions: unexpected status code 404: Comment not found" {
		t.Errorf("expected API error, got %v", err)
	}
}
