package cmd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/testutil"
)

func TestTokenCreateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/my/access_tokens" {
			t.Errorf("expected /test-account/my/access_tokens, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]fizzy.CreateAccessTokenPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}
		got := payload["access_token"]
		if got.Description != "Fizzy CLI" {
			t.Errorf("expected description 'Fizzy CLI', got %s", got.Description)
		}
		if got.Permission != "write" {
			t.Errorf("expected permission 'write', got %s", got.Permission)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fizzy.PersonalAccessToken{
			Token: "secret-value", Description: "Fizzy CLI", Permission: "write",
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := tokenCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--description", "Fizzy CLI", "--permission", "write"})

	if err := handleCreateToken(cmd); err != nil {
		t.Fatalf("handleCreateToken failed: %v", err)
	}
}

func TestTokenCreateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := tokenCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--description", "x", "--permission", "read"})

	err := handleCreateToken(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "creating access token: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestTokenCreateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := tokenCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--description", "x", "--permission", "read"})

	err := handleCreateToken(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
