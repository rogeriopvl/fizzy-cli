package cmd

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/testutil"
)

func TestUserAvatarDeleteCommandSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/test-account/users/u-1/avatar" {
			t.Errorf("expected /test-account/users/u-1/avatar, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userAvatarDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	if err := handleDeleteUserAvatar(cmd, "u-1"); err != nil {
		t.Fatalf("handleDeleteUserAvatar failed: %v", err)
	}
}

func TestUserAvatarDeleteCommandAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := testutil.NewTestClient(server.URL, "", "", "test-token")
	testApp := &app.App{Client: client}

	cmd := userAvatarDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteUserAvatar(cmd, "u-1")
	if err == nil {
		t.Errorf("expected error for API failure")
	}
	if err.Error() != "deleting user avatar: unexpected status code 500: Internal Server Error" {
		t.Errorf("expected API error, got %v", err)
	}
}

func TestUserAvatarDeleteCommandNoClient(t *testing.T) {
	testApp := &app.App{}

	cmd := userAvatarDeleteCmd
	cmd.SetContext(testApp.ToContext(context.Background()))

	err := handleDeleteUserAvatar(cmd, "u-1")
	if err == nil {
		t.Errorf("expected error when client not available")
	}
	if err.Error() != "API client not available" {
		t.Errorf("expected 'client not available' error, got %v", err)
	}
}
