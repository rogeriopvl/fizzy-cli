package cmd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
	"github.com/spf13/cobra"
)

func TestStepUpdateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/cards/123/steps/step-456" {
			t.Errorf("expected /test-account/cards/123/steps/step-456, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		stepPayload := payload["step"]
		if stepPayload["content"] != "Updated step text" {
			t.Errorf("expected content 'Updated step text', got %v", stepPayload["content"])
		}

		w.Header().Set("Content-Type", "application/json")
		response := fizzy.Step{
			ID:        "step-456",
			Content:   "Updated step text",
			Completed: false,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := stepUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--content", "Updated step text",
	})

	if err := handleUpdateStep(cmd, "123", "step-456"); err != nil {
		t.Fatalf("handleUpdateStep failed: %v", err)
	}
}

func TestStepUpdateCommandWithCompleted(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload map[string]map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		stepPayload := payload["step"]
		if stepPayload["completed"] != true {
			t.Errorf("expected completed true, got %v", stepPayload["completed"])
		}

		w.Header().Set("Content-Type", "application/json")
		response := fizzy.Step{
			ID:        "step-456",
			Content:   "Some step",
			Completed: true,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := stepUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--completed",
	})

	if err := handleUpdateStep(cmd, "123", "step-456"); err != nil {
		t.Fatalf("handleUpdateStep failed: %v", err)
	}
}

func TestStepUpdateCommandNoFlags(t *testing.T) {
	testApp := &app.App{Client: testutil.NewTestClient("http://localhost", "", "", "test-token")}

	cmd := stepUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	// Reset flags to ensure no flags are set
	cmd.Flags().Set("content", "")
	cmd.Flags().Set("completed", "false")

	// Create a fresh command to avoid flag state from other tests
	freshCmd := &cobra.Command{}
	freshCmd.Flags().StringP("content", "c", "", "New step content")
	freshCmd.Flags().BoolP("completed", "d", false, "Mark step as completed")
	freshCmd.SetContext(testApp.ToContext(context.Background()))

	err := handleUpdateStep(freshCmd, "123", "step-456")
	if err == nil {
		t.Errorf("expected error when no flags provided")
	}
	if err.Error() != "at least one of --content or --completed must be provided" {
		t.Errorf("expected 'at least one flag' error, got %v", err)
	}
}

func TestStepUpdateCommandInvalidCardNumber(t *testing.T) {
	testApp := &app.App{}

	cmd := stepUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--content", "Updated text",
	})

	err := handleUpdateStep(cmd, "not-a-number", "step-456")
	if err == nil {
		t.Errorf("expected error for invalid card number")
	}
	if err.Error() != "invalid card number: strconv.Atoi: parsing \"not-a-number\": invalid syntax" {
		t.Errorf("expected invalid card number error, got %v", err)
	}
}

func TestStepUpdateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := stepUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--content", "Updated text",
	})

	err := handleUpdateStep(cmd, "123", "step-456")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestStepUpdateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Step not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := stepUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--content", "Updated text",
	})

	err := handleUpdateStep(cmd, "123", "step-456")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "updating step: unexpected status code 404: Step not found" {
		t.Errorf("expected API error, got %v", err)
	}
}
