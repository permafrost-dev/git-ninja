package cmd

import (
	"fmt"

	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/spf13/cobra"
)

var flagCheckout bool = false

var branchLastCmd = &cobra.Command{
	Use:   "branch:last",
	Short: "Work with the last checked out branch",
	Run: func(cmd *cobra.Command, args []string) {
		branchName, _ := helpers.GetLastCheckedoutBranchName()

		if flagCheckout {
			// git prints success and error messages automatically, so we don't need to do it here
			helpers.RunCommandOnStdout("git", "checkout", branchName)
			return
		}

		fmt.Println(branchName)
	},
}

func init() {
	rootCmd.AddCommand(branchLastCmd)

	branchLastCmd.Flags().BoolVarP(&flagCheckout, "checkout", "c", false, "Switch to (checkout) the last checked out branch")
}
