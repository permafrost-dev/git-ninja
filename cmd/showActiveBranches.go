/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/permafrost-dev/git-ninja/app/git"
	"github.com/spf13/cobra"
)

// showActiveBranchesCmd represents the showActiveBranches command
var showActiveBranchesCmd = &cobra.Command{
	Use:   "branch:actives",
	Short: "List active branches, both local and remote",
	Long:  `List active branches, both local and remote.`,
	Run: func(cmd *cobra.Command, args []string) {
		git.ShowActive()
	},
}

func init() {
	rootCmd.AddCommand(showActiveBranchesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showActiveBranchesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showActiveBranchesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
