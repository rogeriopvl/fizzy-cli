package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestAccountListCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/my/identity" {
			t.Errorf("expected /my/identity, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("missing Authorization header")
		}
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := fizzy.GetMyIdentityResponse{
			Accounts: []fizzy.Account{
				{
					ID:        "account-123",
					Name:      "Personal",
					Slug:      "personal",
					CreatedAt: "2025-01-01T00:00:00Z",
				},
				{
					ID:        "account-456",
					Name:      "Work",
					Slug:      "work",
					CreatedAt: "2025-01-02T00:00:00Z",
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := accountListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListAccounts(cmd); err != nil {
		t.Fatalf("handleListAccounts failed: %v", err)
	}
}

func TestAccountListCommandNoAccounts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := fizzy.GetMyIdentityResponse{Accounts: []fizzy.Account{}}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := accountListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleListAccounts(cmd); err != nil {
		t.Fatalf("handleListAccounts failed: %v", err)
	}
}

func TestAccountListCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := accountListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListAccounts(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching accounts: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestAccountListCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := accountListCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleListAccounts(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
