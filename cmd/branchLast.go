package cmd

import (
	"fmt"

	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/spf13/cobra"
)

var flagCheckout bool = false

func switchToBranch(branchName string) {
	err := helpers.RunCommandOnStdout("git", "checkout", branchName)
	if err != nil {
		fmt.Println("Error switching branches:", err)
		return
	}
}

var branchLastCmd = &cobra.Command{
	Use:   "branch:last",
	Short: "Work with the last checked out branch",
	Run: func(cmd *cobra.Command, args []string) {
		branchName, _ := helpers.GetLastCheckedoutBranchName()

		if flagCheckout {
			switchToBranch(branchName)
			return
		}

		fmt.Println(branchName)
	},
}

func init() {
	rootCmd.AddCommand(branchLastCmd)

	branchLastCmd.Flags().BoolVarP(&flagCheckout, "switch", "s", false, "Switch to the last checked out branch")
}
