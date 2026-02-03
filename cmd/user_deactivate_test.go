package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestUserDeactivateCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/user-123" {
			t.Errorf("expected /users/user-123, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userDeactivateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleDeactivateUser(cmd, "user-123"); err != nil {
		t.Fatalf("handleDeactivateUser failed: %v", err)
	}
}

func TestUserDeactivateCommandNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userDeactivateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeactivateUser(cmd, "nonexistent-user")
	if err == nil {
		t.Errorf("expected error for user not found")
	}
	if err.Error() != "deactivating user: unexpected status code 404: User not found" {
		t.Errorf("expected user not found error, got %v", err)
	}
}

func TestUserDeactivateCommandForbidden(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("You don't have permission to deactivate this user"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userDeactivateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeactivateUser(cmd, "user-123")
	if err == nil {
		t.Errorf("expected error for forbidden access")
	}
	if err.Error() != "deactivating user: unexpected status code 403: You don't have permission to deactivate this user" {
		t.Errorf("expected permission error, got %v", err)
	}
}

func TestUserDeactivateCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userDeactivateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeactivateUser(cmd, "user-123")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "deactivating user: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestUserDeactivateCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := userDeactivateCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeactivateUser(cmd, "user-123")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
