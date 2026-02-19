// Package app
package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rogeriopvl/fizzy/internal/config"
	fizzy "github.com/rogeriopvl/fizzy-go"
)

type App struct {
	Client *fizzy.Client
	Config *config.Config
}

func New(version string) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}

	token, isSet := os.LookupEnv("FIZZY_ACCESS_TOKEN")
	if !isSet || token == "" {
		return &App{Config: cfg}, nil // No token set, app will handle gracefully
	}

	opts := []fizzy.ClientOption{
		fizzy.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
	}

	if cfg.SelectedBoard != "" {
		opts = append(opts, fizzy.WithBoard(cfg.SelectedBoard))
	}

	client, err := fizzy.NewClient(cfg.SelectedAccount, token, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating API client: %w", err)
	}

	return &App{Client: client, Config: cfg}, nil
}

// contextKey is a type for context keys to avoid collisions.
type contextKey string

const appContextKey contextKey = "app"

func FromContext(ctx context.Context) *App {
	app, _ := ctx.Value(appContextKey).(*App)
	return app
}

func (a *App) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, appContextKey, a)
}
