package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestStepDeleteCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/123/steps/step-456" {
			t.Errorf("expected /test-account/cards/123/steps/step-456, got %s", r.URL.Path)
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

	cmd := stepDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleDeleteStep(cmd, "123", "step-456"); err != nil {
		t.Fatalf("handleDeleteStep failed: %v", err)
	}
}

func TestStepDeleteCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := stepDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteStep(cmd, "not-a-number", "step-456")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}

func TestStepDeleteCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := stepDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteStep(cmd, "123", "step-456")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestStepDeleteCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Step not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := stepDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteStep(cmd, "123", "step-456")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "deleting step: unexpected status code 404: Step not found" {
		t.Errorf("expected API error, got %v", err)
	}
}
