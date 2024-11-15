package cmd

import (
	"fmt"

	"github.com/permafrost-dev/git-ninja/app/git"
	"github.com/spf13/cobra"
)

var showActiveBranchesCmd = &cobra.Command{
	Use:    "branch:actives",
	Hidden: true,
	Short:  "List active branches, both local and remote",
	Long:   `List active branches, both local and remote.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("showActiveBranches called")

		git.ShowActive()
	},
}

func init() {
	rootCmd.AddCommand(showActiveBranchesCmd)
}
