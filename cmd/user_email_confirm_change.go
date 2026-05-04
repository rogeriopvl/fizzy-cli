package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var userEmailConfirmChangeCmd = &cobra.Command{
	Use:   "confirm-change <user_id>",
	Short: "Confirm a user email address change",
	Long:  `Confirm a previously-requested email address change using the token sent to the new address.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleConfirmUserEmailChange(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleConfirmUserEmailChange(cmd *cobra.Command, userID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	token, _ := cmd.Flags().GetString("token")

	if err := a.Client.ConfirmUserEmailChange(context.Background(), userID, token); err != nil {
		return fmt.Errorf("confirming email change: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Email change confirmed for user %s\n", userID)
	return nil
}

func init() {
	userEmailConfirmChangeCmd.Flags().StringP("token", "t", "", "Confirmation token (required)")
	userEmailConfirmChangeCmd.MarkFlagRequired("token")

	userEmailCmd.AddCommand(userEmailConfirmChangeCmd)
}
