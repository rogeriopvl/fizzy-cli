package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestWhoamiCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/my/identity" {
			t.Errorf("expected /my/identity, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		response := api.GetMyIdentityResponse{
			Accounts: []api.Account{
				{
					ID:   "account-123",
					Name: "My Company",
					Slug: "/123456",
					User: api.User{
						ID:    "user-1",
						Name:  "John Doe",
						Email: "john@example.com",
						Role:  "owner",
					},
				},
				{
					ID:   "account-456",
					Name: "Another Co",
					Slug: "/789012",
					User: api.User{
						ID:    "user-1",
						Name:  "John Doe",
						Email: "john@example.com",
						Role:  "member",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := whoamiCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleWhoami(cmd); err != nil {
		t.Fatalf("handleWhoami failed: %v", err)
	}
}

func TestWhoamiCommandNoAccounts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := api.GetMyIdentityResponse{Accounts: []api.Account{}}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := whoamiCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleWhoami(cmd); err != nil {
		t.Fatalf("handleWhoami failed: %v", err)
	}
}

func TestWhoamiCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := whoamiCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleWhoami(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
}

func TestWhoamiCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := whoamiCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleWhoami(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
