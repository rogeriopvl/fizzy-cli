package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/testutil"
)

func TestCardImageDeleteCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/123/image" {
			t.Errorf("expected /test-account/cards/123/image, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardImageDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleDeleteCardImage(cmd, "123"); err != nil {
		t.Fatalf("handleDeleteCardImage failed: %v", err)
	}
}

func TestCardImageDeleteCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := cardImageDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteCardImage(cmd, "not-a-number")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
}

func TestCardImageDeleteCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardImageDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteCardImage(cmd, "123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "deleting card image: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestCardImageDeleteCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := cardImageDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteCardImage(cmd, "123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
