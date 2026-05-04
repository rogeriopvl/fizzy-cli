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

func TestExportAccountCreateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/account/exports" {
			t.Errorf("expected /test-account/account/exports, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fizzy.Export{
			ID: "exp-1", Status: "pending", CreatedAt: "2026-04-02T12:34:56Z",
		})
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := exportAccountCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleCreateAccountExport(cmd); err != nil {
		t.Fatalf("handleCreateAccountExport failed: %v", err)
	}
}

func TestExportAccountCreateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := exportAccountCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleCreateAccountExport(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "creating account export: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestExportAccountCreateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := exportAccountCreateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleCreateAccountExport(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
