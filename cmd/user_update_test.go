package cmd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
	"github.com/spf13/cobra"
)

func TestUserUpdateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/user-123" {
			t.Errorf("expected /users/user-123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]api.UpdateUserPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("failed to unmarshal request body: %v", err)
		}

		userPayload := payload["user"]
		if userPayload.Name != "Updated Name" {
			t.Errorf("expected name 'Updated Name', got %s", userPayload.Name)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Updated Name"})

	if err := handleUpdateUser(cmd, "user-123"); err != nil {
		t.Fatalf("handleUpdateUser failed: %v", err)
	}
}

func TestUserUpdateCommandWithAllFlags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var payload map[string]api.UpdateUserPayload
		json.Unmarshal(body, &payload)

		userPayload := payload["user"]
		if userPayload.Name != "Updated Name" {
			t.Errorf("expected name 'Updated Name', got %s", userPayload.Name)
		}
		if userPayload.Avatar != "https://example.com/avatar.png" {
			t.Errorf("expected avatar 'https://example.com/avatar.png', got %s", userPayload.Avatar)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{
		"--name", "Updated Name",
		"--avatar", "https://example.com/avatar.png",
	})

	if err := handleUpdateUser(cmd, "user-123"); err != nil {
		t.Fatalf("handleUpdateUser failed: %v", err)
	}
}

func TestUserUpdateCommandNoFlags(t *testing.T) {
	client := testutil.NewTestClient("http://localhost", "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := &cobra.Command{
		Use:  "update <user_id>",
		Args: cobra.ExactArgs(1),
	}
	cmd.Flags().String("name", "", "User name")
	cmd.Flags().String("avatar", "", "Avatar URL or file path")

	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleUpdateUser(cmd, "user-123")
	if err == nil {
		t.Errorf("expected error when no flags provided")
	}
	if err.Error() != "at least one flag must be provided (--name or --avatar)" {
		t.Errorf("expected flag requirement error, got %v", err)
	}
}

func TestUserUpdateCommandNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Updated Name"})

	err := handleUpdateUser(cmd, "nonexistent-user")
	if err == nil {
		t.Errorf("expected error for user not found")
	}
	if err.Error() != "updating user: unexpected status code 404: User not found" {
		t.Errorf("expected user not found error, got %v", err)
	}
}

func TestUserUpdateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Updated Name"})

	err := handleUpdateUser(cmd, "user-123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "updating user: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestUserUpdateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := userUpdateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))
	cmd.ParseFlags([]string{"--name", "Updated Name"})

	err := handleUpdateUser(cmd, "user-123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
