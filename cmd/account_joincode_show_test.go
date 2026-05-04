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

func TestAccountJoincodeShowCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/account/join_code" {
			t.Errorf("expected /test-account/account/join_code, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fizzy.JoinCode{
			Code: "abc123", UsageCount: 3, UsageLimit: 10, URL: "https://example/join/abc123", Active: true,
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := accountJoincodeShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleShowJoinCode(cmd); err != nil {
		t.Fatalf("handleShowJoinCode failed: %v", err)
	}
}

func TestAccountJoincodeShowCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := accountJoincodeShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowJoinCode(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching join code: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestAccountJoincodeShowCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := accountJoincodeShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowJoinCode(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
