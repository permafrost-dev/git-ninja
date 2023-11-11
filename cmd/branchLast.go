package cmd

import (
	"fmt"

	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "branch:last",
		Short: "Show the last checked out branch name",
		Run: func(cmd *cobra.Command, args []string) {
			branchName, _ := helpers.GetLastCheckedoutBranchName()
			fmt.Println(branchName)
		},
	})
}
