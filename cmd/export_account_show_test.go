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

func TestExportAccountShowCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/account/exports/exp-1" {
			t.Errorf("expected /test-account/account/exports/exp-1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fizzy.Export{
			ID: "exp-1", Status: "completed", CreatedAt: "2026-04-02T12:34:56Z", DownloadURL: "https://example/x",
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := exportAccountShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleShowAccountExport(cmd, "exp-1"); err != nil {
		t.Fatalf("handleShowAccountExport failed: %v", err)
	}
}

func TestExportAccountShowCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := exportAccountShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowAccountExport(cmd, "exp-1")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "fetching account export: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestExportAccountShowCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := exportAccountShowCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleShowAccountExport(cmd, "exp-1")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
