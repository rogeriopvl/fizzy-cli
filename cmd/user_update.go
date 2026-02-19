package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/spf13/cobra"
)

var userUpdateCmd = &cobra.Command{
	Use:   "update <user_id>",
	Short: "Update a user",
	Long: `Update user settings such as name and avatar.

Avatar must be provided as a URL (e.g., https://example.com/avatar.jpg).

Example:
  fizzy user update user-123 --name "John Doe"
  fizzy user update user-123 --avatar https://example.com/avatar.jpg`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUpdateUser(cmd, args[0]); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUpdateUser(cmd *cobra.Command, userID string) error {
	if !cmd.Flags().Changed("name") && !cmd.Flags().Changed("avatar") {
		return fmt.Errorf("at least one flag must be provided (--name or --avatar)")
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	payload := fizzy.UpdateUserPayload{}

	if cmd.Flags().Changed("name") {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return fmt.Errorf("invalid name flag: %w", err)
		}
		payload.Name = name
	}
	if cmd.Flags().Changed("avatar") {
		avatar, err := cmd.Flags().GetString("avatar")
		if err != nil {
			return fmt.Errorf("invalid avatar flag: %w", err)
		}
		payload.Avatar = avatar
	}

	err := a.Client.UpdateUser(context.Background(), userID, payload)
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ User '%s' updated successfully\n", userID)
	return nil
}

func init() {
	userUpdateCmd.Flags().StringP("name", "n", "", "User name")
	userUpdateCmd.Flags().String("avatar", "", "Avatar URL (e.g., https://example.com/avatar.jpg)")

	userCmd.AddCommand(userUpdateCmd)
}
