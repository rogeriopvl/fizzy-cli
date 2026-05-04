package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var userEmailRequestChangeCmd = &cobra.Command{
	Use:   "request-change <user_id>",
	Short: "Request a user email address change",
	Long:  `Request an email address change for a user. A confirmation token will be sent to the new address.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleRequestUserEmailChange(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleRequestUserEmailChange(cmd *cobra.Command, userID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	email, _ := cmd.Flags().GetString("email")
	payload := fizzy.RequestEmailChangePayload{EmailAddress: email}

	if err := a.Client.RequestUserEmailChange(context.Background(), userID, payload); err != nil {
		return fmt.Errorf("requesting email change: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Email change requested for user %s; check %s for a confirmation token\n", userID, email)
	return nil
}

func init() {
	userEmailRequestChangeCmd.Flags().StringP("email", "e", "", "New email address (required)")
	userEmailRequestChangeCmd.MarkFlagRequired("email")

	userEmailCmd.AddCommand(userEmailRequestChangeCmd)
}
