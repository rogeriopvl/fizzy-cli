package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestCardReactionDeleteCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/123/reactions/reaction-456" {
			t.Errorf("expected /test-account/cards/123/reactions/reaction-456, got %s", r.URL.Path)
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

	cmd := cardReactionDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleDeleteCardReaction(cmd, "123", "reaction-456"); err != nil {
		t.Fatalf("handleDeleteCardReaction failed: %v", err)
	}
}

func TestCardReactionDeleteCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := cardReactionDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteCardReaction(cmd, "not-a-number", "reaction-456")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}

func TestCardReactionDeleteCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := cardReactionDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteCardReaction(cmd, "123", "reaction-456")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestCardReactionDeleteCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Reaction not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardReactionDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteCardReaction(cmd, "123", "reaction-456")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "deleting reaction: unexpected status code 404: Reaction not found" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestCardReactionDeleteCommandDifferentReactionID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/456/reactions/reaction-789" {
			t.Errorf("expected /test-account/cards/456/reactions/reaction-789, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardReactionDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleDeleteCardReaction(cmd, "456", "reaction-789"); err != nil {
		t.Fatalf("handleDeleteCardReaction failed: %v", err)
	}
}
