package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy-cli/internal/api"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/config"
	"github.com/rogeriopvl/fizzy-cli/internal/testutil"
	"github.com/spf13/cobra"
)

func newUseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use",
		Short: "Set the active board or account",
		Long:  `Set the active board or account to use for subsequent commands`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := handleUse(cmd); err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
			}
		},
	}
	cmd.Flags().String("board", "", "Board name to use")
	cmd.Flags().String("account", "", "Account slug to use")
	return cmd
}

func TestUseCommandSetBoard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/boards" {
			t.Errorf("expected /boards, got %s", r.URL.Path)
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
		response := []api.Board{
			{
				ID:   "board-123",
				Name: "My Project",
			},
			{
				ID:   "board-456",
				Name: "Other Board",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	cfg := &config.Config{}
	testApp := &app.App{
		Client: client,
		Config: cfg,
	}

	cmd := newUseCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board", "My Project"})

	if err := handleUse(cmd); err != nil {
		t.Fatalf("handleUse failed: %v", err)
	}

	savedCfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if savedCfg.SelectedBoard != "board-123" {
		t.Errorf("expected SelectedBoard=board-123, got %s", savedCfg.SelectedBoard)
	}
}

func TestUseCommandSetAccount(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cmd := newUseCmd()
	cmd.ParseFlags([]string{"--account", "my-company"})

	cfg := &config.Config{}
	testApp := &app.App{Config: cfg}
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleUse(cmd); err != nil {
		t.Fatalf("handleUse failed: %v", err)
	}

	savedCfg, err := config.Load()
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if savedCfg.SelectedAccount != "my-company" {
		t.Errorf("expected SelectedAccount=my-company, got %s", savedCfg.SelectedAccount)
	}
}

func TestUseCommandBoardNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.Header().Set("Content-Type", "application/json")
		response := []api.Board{
			{ID: "board-123", Name: "Existing Board"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{
		Client: client,
		Config: &config.Config{},
	}

	cmd := newUseCmd()
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--board", "Nonexistent Board"})

	err := handleUse(cmd)
	if err == nil {
		t.Errorf("expected error for nonexistent board")
	}
	if err.Error() != "board 'Nonexistent Board' not found" {
		t.Errorf("expected 'not found' error, got %v", err)
	}
}

func TestUseCommandBothFlagsError(t *testing.T) {
	cmd := newUseCmd()
	cmd.ParseFlags([]string{"--board", "Board", "--account", "Account"})
	cmd.SetContext(context.Background())

	err := handleUse(cmd)
	if err == nil {
		t.Errorf("expected error when both flags provided")
	}
	if err.Error() != "cannot specify both --board and --account" {
		t.Errorf("expected 'both flags' error, got %v", err)
	}
}

func TestUseCommandNoFlagsError(t *testing.T) {
	cmd := newUseCmd()
	cmd.SetContext(context.Background())
	cmd.ParseFlags([]string{})

	err := handleUse(cmd)
	if err == nil {
		t.Errorf("expected error when no flags provided")
	}
	if err.Error() != "must specify either --board or --account" {
		t.Errorf("expected 'no flags' error, got %v", err)
	}
}
