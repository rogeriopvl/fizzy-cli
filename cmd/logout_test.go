package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/testutil"
)

func TestLogoutCommand(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/session" {
			t.Errorf("expected /session, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("missing Authorization header")
		}
		if auth != "Bearer test-token" {
			t.Errorf("expected Bearer test-token, got %s", auth)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := logoutCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleLogout(cmd); err != nil {
		t.Fatalf("handleLogout failed: %v", err)
	}
}

func TestLogoutCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := logoutCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleLogout(cmd)
	if err == nil {
		t.Errorf("expected error for API failure")
	}
}

func TestLogoutCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := logoutCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleLogout(cmd)
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}

func TestLogoutCommandUnauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "invalid-token")
	testApp := &app.App{Client: client}

	cmd := logoutCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleLogout(cmd)
	if err == nil {
		t.Errorf("expected error for unauthorized request")
	}
}
