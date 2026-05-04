package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var accountJoincodeUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the account join code",
	Long:  `Update the join code's usage limit. Requires admin role.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleUpdateJoinCode(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleUpdateJoinCode(cmd *cobra.Command) error {
	if !cmd.Flags().Changed("usage-limit") {
		return fmt.Errorf("--usage-limit is required")
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	limit, _ := cmd.Flags().GetInt("usage-limit")
	payload := fizzy.UpdateJoinCodePayload{UsageLimit: limit}

	if err := a.Client.UpdateAccountJoinCode(context.Background(), payload); err != nil {
		return fmt.Errorf("updating join code: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Join code usage limit set to %d\n", limit)
	return nil
}

func init() {
	accountJoincodeUpdateCmd.Flags().Int("usage-limit", 0, "Maximum number of times the join code can be used (required)")
	accountJoincodeCmd.AddCommand(accountJoincodeUpdateCmd)
}
