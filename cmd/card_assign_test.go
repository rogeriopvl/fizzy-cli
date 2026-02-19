package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/config"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestCardAssignCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/123/assignments" {
			t.Errorf("expected /test-account/cards/123/assignments, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
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

	cmd := cardAssignCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleAssignCard(cmd, "123", "user-id-123"); err != nil {
		t.Fatalf("handleAssignCard failed: %v", err)
	}
}

func TestCardAssignCommandWithMe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/123/assignments" {
			t.Errorf("expected /test-account/cards/123/assignments, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	cfg := &config.Config{CurrentUserID: "my-user-id"}
	testApp := &app.App{Client: client, Config: cfg}

	cmd := cardAssignCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleAssignCard(cmd, "123", "me"); err != nil {
		t.Fatalf("handleAssignCard with 'me' failed: %v", err)
	}
}

func TestCardAssignCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Card not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardAssignCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleAssignCard(cmd, "999", "user-id-123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "assigning card: unexpected status code 404: Card not found" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestCardAssignCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := cardAssignCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleAssignCard(cmd, "not-a-number", "user-id-123")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}

func TestCardAssignCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := cardAssignCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleAssignCard(cmd, "123", "user-id-123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestCardAssignCommandMeWithoutUserID(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{Client: client, Config: &config.Config{}}

	cmd := cardAssignCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleAssignCard(cmd, "123", "me")
	if err == nil {
		t.Errorf("expected error when using 'me' without current user ID")
	}
	if err.Error() != "current user ID not available, please run 'fizzy login' first" {
		t.Errorf("expected 'current user ID not available' error, got %v", err)
	}
}
