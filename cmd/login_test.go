package cmd

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/config"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestLoginCommand(t *testing.T) {
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
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("expected Accept: application/json, got %s", r.Header.Get("Accept"))
		}

		w.Header().Set("Content-Type", "application/json")
		response := fizzy.GetMyIdentityResponse{
			Accounts: []fizzy.Account{
				{
					ID:   "acc-123",
					Name: "Test Account",
					Slug: "test-account",
					User: fizzy.User{
						ID:    "user-123",
						Email: "test@example.com",
						Name:  "Test User",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	configPath := tmpDir + "/.config/fizzy-cli/config.json"

	t.Setenv("FIZZY_ACCESS_TOKEN", "test-token")
	t.Setenv("HOME", tmpDir)

	client := testutil.NewTestClient(server.URL, "", "", "test-token")

	cfg := &config.Config{}
	testApp := &app.App{
		Client: client,
		Config: cfg,
	}

	cmd := loginCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleLogin(cmd); err != nil {
		t.Fatalf("handleLogin failed: %v", err)
	}

	if cfg.SelectedAccount != "test-account" {
		t.Errorf("expected SelectedAccount=test-account, got %s", cfg.SelectedAccount)
	}

	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("config file not created at %s: %v", configPath, err)
	}
}

func TestLoginCommandWithoutToken(t *testing.T) {
	t.Setenv("FIZZY_ACCESS_TOKEN", "")

	cmd := loginCmd
	err := handleLogin(cmd)

	if err != nil {
		t.Errorf("expected no error when token is missing, got %v", err)
	}
}
