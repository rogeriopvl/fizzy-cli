package cmd

import (
	"context"
	"fmt"

	"github.com/rogeriopvl/fizzy/internal/api"
	"github.com/rogeriopvl/fizzy/internal/app"
	"github.com/rogeriopvl/fizzy/internal/ui"
	"github.com/spf13/cobra"
)

var cardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cards",
	Long: `Retrieve and display cards from Fizzy with optional filters.

Filter options:
  --tag <id>          Filter by tag ID (can be used multiple times)
  --assignee <id>     Filter by assignee user ID (can be used multiple times)
  --creator <id>      Filter by creator user ID (can be used multiple times)
  --closer <id>       Filter by user who closed the card (can be used multiple times)
  --card <id>         Filter to specific card ID (can be used multiple times)
  --indexed-by        Filter by status: all, closed, not_now, stalled, postponing_soon, golden
  --sorted-by         Sort order: latest, newest, oldest
  --unassigned        Show only unassigned cards
  --created-in        Filter by creation date: today, yesterday, thisweek, lastweek, thismonth, lastmonth, thisyear, lastyear
  --closed-in         Filter by closure date: today, yesterday, thisweek, lastweek, thismonth, lastmonth, thisyear, lastyear
  --search            Search terms (can be used multiple times)`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleListCards(cmd); err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error: %v\n", err)
		}
	},
}

func handleListCards(cmd *cobra.Command) error {
	a := app.FromContext(cmd.Context())
	if a == nil || a.Client == nil {
		return fmt.Errorf("API client not available")
	}

	if a.Config.SelectedBoard == "" {
		return fmt.Errorf("no board selected")
	}

	filters := api.CardFilters{
		BoardIDs: []string{a.Config.SelectedBoard},
	}

	if tags, _ := cmd.Flags().GetStringSlice("tag"); len(tags) > 0 {
		filters.TagIDs = tags
	}
	if assignees, _ := cmd.Flags().GetStringSlice("assignee"); len(assignees) > 0 {
		filters.AssigneeIDs = assignees
	}
	if creators, _ := cmd.Flags().GetStringSlice("creator"); len(creators) > 0 {
		filters.CreatorIDs = creators
	}
	if closers, _ := cmd.Flags().GetStringSlice("closer"); len(closers) > 0 {
		filters.CloserIDs = closers
	}
	if cardIDs, _ := cmd.Flags().GetStringSlice("card"); len(cardIDs) > 0 {
		filters.CardIDs = cardIDs
	}
	if searches, _ := cmd.Flags().GetStringSlice("search"); len(searches) > 0 {
		filters.Terms = searches
	}

	if indexedBy, _ := cmd.Flags().GetString("indexed-by"); indexedBy != "" {
		filters.IndexedBy = indexedBy
	}
	if sortedBy, _ := cmd.Flags().GetString("sorted-by"); sortedBy != "" {
		filters.SortedBy = sortedBy
	}
	if createdIn, _ := cmd.Flags().GetString("created-in"); createdIn != "" {
		filters.CreationStatus = createdIn
	}
	if closedIn, _ := cmd.Flags().GetString("closed-in"); closedIn != "" {
		filters.ClosureStatus = closedIn
	}

	if unassigned, _ := cmd.Flags().GetBool("unassigned"); unassigned {
		filters.AssignmentStatus = "unassigned"
	}

	cards, err := a.Client.GetCards(context.Background(), filters)
	if err != nil {
		return fmt.Errorf("fetching cards: %w", err)
	}

	if len(cards) == 0 {
		fmt.Println("No cards found")
		return nil
	}

	return ui.DisplayCards(cards)
}

func init() {
	cardListCmd.Flags().StringSliceP("tag", "t", []string{}, "Filter by tag ID (can be used multiple times)")
	cardListCmd.Flags().StringSliceP("assignee", "a", []string{}, "Filter by assignee user ID (can be used multiple times)")
	cardListCmd.Flags().StringSlice("creator", []string{}, "Filter by creator user ID (can be used multiple times)")
	cardListCmd.Flags().StringSlice("closer", []string{}, "Filter by closer user ID (can be used multiple times)")
	cardListCmd.Flags().StringSlice("card", []string{}, "Filter to specific card ID (can be used multiple times)")
	cardListCmd.Flags().String("indexed-by", "", "Filter by status: all, closed, not_now, stalled, postponing_soon, golden")
	cardListCmd.Flags().String("sorted-by", "", "Sort order: latest, newest, oldest")
	cardListCmd.Flags().BoolP("unassigned", "u", false, "Show only unassigned cards")
	cardListCmd.Flags().String("created-in", "", "Filter by creation date")
	cardListCmd.Flags().String("closed-in", "", "Filter by closure date")
	cardListCmd.Flags().StringSliceP("search", "s", []string{}, "Search terms (can be used multiple times)")

	cardCmd.AddCommand(cardListCmd)
}
