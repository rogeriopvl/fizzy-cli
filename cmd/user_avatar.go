package cmd

import "github.com/spf13/cobra"

var userAvatarCmd = &cobra.Command{
	Use:   "avatar",
	Short: "Manage user avatars",
	Long:  `Manage avatar images for users`,
}

func init() {
	userCmd.AddCommand(userAvatarCmd)
}
