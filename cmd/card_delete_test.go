package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestCardDeleteCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cards/1" {
			t.Errorf("expected /cards/1, got %s", r.URL.Path)
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

	cmd := cardDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleDeleteCard(cmd, "1"); err != nil {
		t.Fatalf("handleDeleteCard failed: %v", err)
	}
}

func TestCardDeleteCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Card not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteCard(cmd, "999")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "deleting card: unexpected status code 404: Card not found" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestCardDeleteCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := cardDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteCard(cmd, "1")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestCardDeleteCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := cardDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteCard(cmd, "not-a-number")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}
