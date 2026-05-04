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

func TestExportUserShowCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/users/user-1/data_exports/exp-9" {
			t.Errorf("expected /test-account/users/user-1/data_exports/exp-9, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fizzy.Export{
			ID: "exp-9", Status: "completed", CreatedAt: "2026-04-02T12:34:56Z", DownloadURL: "https://example/y",
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := exportUserShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleShowUserDataExport(cmd, "user-1", "exp-9"); err != nil {
		t.Fatalf("handleShowUserDataExport failed: %v", err)
	}
}

func TestExportUserShowCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := exportUserShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowUserDataExport(cmd, "user-1", "exp-9")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching user data export: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestExportUserShowCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := exportUserShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowUserDataExport(cmd, "user-1", "exp-9")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
