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
	"github.com/spf13/cobra"
)

func TestAccountJoincodeUpdateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/account/join_code" {
			t.Errorf("expected /test-account/account/join_code, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]fizzy.UpdateJoinCodePayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}
		if payload["account_join_code"].UsageLimit != 25 {
			t.Errorf("expected usage_limit=25, got %d", payload["account_join_code"].UsageLimit)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := accountJoincodeUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--usage-limit", "25"})

	if err := handleUpdateJoinCode(cmd); err != nil {
		t.Fatalf("handleUpdateJoinCode failed: %v", err)
	}
}

func TestAccountJoincodeUpdateCommandMissingFlag(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := &cobra.Command{}
	cmd.Flags().Int("usage-limit", 0, "")
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleUpdateJoinCode(cmd)
	if err == nil {
		t.Errorf("expected error when --usage-limit not provided")
	}
	if err.Error() != "--usage-limit is required" {
		t.Errorf("expected required-flag error, got %v", err)
	}
}

func TestAccountJoincodeUpdateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := accountJoincodeUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--usage-limit", "25"})

	err := handleUpdateJoinCode(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "updating join code: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestAccountJoincodeUpdateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := accountJoincodeUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--usage-limit", "25"})

	err := handleUpdateJoinCode(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
