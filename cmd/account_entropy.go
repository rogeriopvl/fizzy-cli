package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/spf13/cobra"
)

var accountEntropyCmd = &cobra.Command{
	Use:   "entropy",
	Short: "Update the account auto-postpone period",
	Long:  `Update the account-level default auto-postpone period (in days). Requires admin role.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleAccountEntropy(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleAccountEntropy(cmd *cobra.Command) error {
	if !cmd.Flags().Changed("auto-postpone-days") {
		return fmt.Errorf("--auto-postpone-days is required")
	}

	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	days, _ := cmd.Flags().GetInt("auto-postpone-days")
	payload := fizzy.EntropyPayload{AutoPostponePeriodInDays: days}

	account, err := a.Client.UpdateAccountEntropy(context.Background(), payload)
	if err != nil {
		return fmt.Errorf("updating account entropy: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "✓ Account auto-postpone period set to %d days\n", account.AutoPostponePeriodInDays)
	return nil
}

func init() {
	accountEntropyCmd.Flags().Int("auto-postpone-days", 0, "Auto-postpone period in days (required)")
	accountCmd.AddCommand(accountEntropyCmd)
}
