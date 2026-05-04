package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var userAvatarDeleteCmd = &cobra.Command{
	Use:   "delete <user_id>",
	Short: "Delete a user's avatar",
	Long:  `Remove the avatar image for a user`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleDeleteUserAvatar(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleDeleteUserAvatar(cmd *cobra.Command, userID string) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	if err := a.Client.DeleteUserAvatar(context.Background(), userID); err != nil {
		return fmt.Errorf("deleting user avatar: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Avatar removed from user %s\n", userID)
	return nil
}

func init() {
	userAvatarCmd.AddCommand(userAvatarDeleteCmd)
}
