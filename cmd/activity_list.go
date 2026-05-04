package cmd

import (
	"context"
	"fmt"

	fizzy "github.com/rogeriopvl/fizzy-go"
	"github.com/rogeriopvl/fizzy-cli/internal/app"
	"github.com/rogeriopvl/fizzy-cli/internal/ui"
	"github.com/spf13/cobra"
)

var activityListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent activities",
	Long:  `Retrieve and display the account's activity feed, newest first`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListActivities(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListActivities(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	filters := &fizzy.ActivityFilters{}
	if creators, _ := cmd.Flags().GetStringSlice("creator"); len(creators) > 0 {
		filters.CreatorIDs = creators
	}
	if boards, _ := cmd.Flags().GetStringSlice("board"); len(boards) > 0 {
		filters.BoardIDs = boards
	}
	if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
		filters.Limit = limit
	}

	activities, err := a.Client.GetActivities(context.Background(), filters)
	if err != nil {
		return fmt.Errorf("fetching activities: %w", err)
	}

	if len(activities) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No activities found")
		return nil
	}

	return ui.DisplayActivities(cmd.OutOrStdout(), activities)
}

func init() {
	activityListCmd.Flags().StringSlice("creator", nil, "Filter by creator user ID (can be used multiple times)")
	activityListCmd.Flags().StringSlice("board", nil, "Filter by board ID (can be used multiple times)")
	activityListCmd.Flags().IntP("limit", "l", 0, "Maximum number of activities to return (0 = no limit)")

	activityCmd.AddCommand(activityListCmd)
}
