package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/testutil"
)

func TestStepShowCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/123/steps/step-1" {
			t.Errorf("expected /test-account/cards/123/steps/step-1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fizzy.Step{ID: "step-1", Content: "Buy milk", Completed: true})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := stepShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleShowStep(cmd, "123", "step-1"); err != nil {
		t.Fatalf("handleShowStep failed: %v", err)
	}
}

func TestStepShowCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := stepShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowStep(cmd, "not-a-number", "step-1")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
}

func TestStepShowCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := stepShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowStep(cmd, "123", "step-1")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching step: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestStepShowCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := stepShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowStep(cmd, "123", "step-1")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
