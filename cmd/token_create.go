package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var tokenCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a personal access token",
	Long:  `Create a new personal access token. The token value is shown once and cannot be retrieved again.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleCreateToken(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleCreateToken(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	description, _ := cmd.Flags().GetString("description")
	permission, _ := cmd.Flags().GetString("permission")

	payload := fizzy.CreateAccessTokenPayload{
		Description: description,
		Permission:  permission,
	}

	token, err := a.Client.CreateAccessToken(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("creating access token: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Token created (description: %s, permission: %s)\n", token.Description, token.Permission)
	fmt.Fprintf(cmd.OutOrStdout(), "Token: %s\n", token.Token)
	fmt.Fprintf(cmd.OutOrStdout(), "Save this value now — it cannot be retrieved again.\n")
	return nil
}

func init() {
	tokenCreateCmd.Flags().StringP("description", "d", "", "Token description (required)")
	tokenCreateCmd.MarkFlagRequired("description")
	tokenCreateCmd.Flags().StringP("permission", "p", "", "Permission: read or write (required)")
	tokenCreateCmd.MarkFlagRequired("permission")

	tokenCmd.AddCommand(tokenCreateCmd)
}
