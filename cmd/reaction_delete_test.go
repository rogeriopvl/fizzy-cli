package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestReactionDeleteCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cards/123/comments/comment-456/reactions/reaction-789" {
			t.Errorf("expected /cards/123/comments/comment-456/reactions/reaction-789, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := reactionDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleDeleteReaction(cmd, "123", "comment-456", "reaction-789"); err != nil {
		t.Fatalf("handleDeleteReaction failed: %v", err)
	}
}

func TestReactionDeleteCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := reactionDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteReaction(cmd, "not-a-number", "comment-456", "reaction-789")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}

func TestReactionDeleteCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := reactionDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteReaction(cmd, "123", "comment-456", "reaction-789")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestReactionDeleteCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Reaction not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := reactionDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteReaction(cmd, "123", "comment-456", "reaction-789")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "deleting reaction: unexpected status code 404: Reaction not found" {
		t.Errorf("expected API error, got %v", err)
	}
}
