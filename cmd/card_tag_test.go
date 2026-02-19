package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestCardTagCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/123/taggings" {
			t.Errorf("expected /test-account/cards/123/taggings, got %s", r.URL.Path)
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

	cmd := cardTagCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleTagCard(cmd, "123", "bug"); err != nil {
		t.Fatalf("handleTagCard failed: %v", err)
	}
}

func TestCardTagCommandWithHashPrefix(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/456/taggings" {
			t.Errorf("expected /test-account/cards/456/taggings, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardTagCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleTagCard(cmd, "456", "#feature"); err != nil {
		t.Fatalf("handleTagCard with # prefix failed: %v", err)
	}
}

func TestCardTagCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Card not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := cardTagCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleTagCard(cmd, "999", "bug")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "toggling tag on card: unexpected status code 404: Card not found" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestCardTagCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := cardTagCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleTagCard(cmd, "not-a-number", "bug")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}

func TestCardTagCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := cardTagCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleTagCard(cmd, "123", "bug")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
