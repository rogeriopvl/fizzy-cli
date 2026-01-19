package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestCardUngoldenCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cards/123/goldness" {
			t.Errorf("expected /cards/123/goldness, got %s", r.URL.Path)
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

	cmd := cardUngoldenCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleUngoldenCard(cmd, "123"); err != nil {
		t.Fatalf("handleUngoldenCard failed: %v", err)
	}
}

func TestCardUngoldenCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := cardUngoldenCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleUngoldenCard(cmd, "not-a-number")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}

func TestCardUngoldenCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := cardUngoldenCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleUngoldenCard(cmd, "123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestCardUngoldenCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardUngoldenCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleUngoldenCard(cmd, "123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "removing golden status: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}
