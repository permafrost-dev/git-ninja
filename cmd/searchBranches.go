package cmd

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/permafrost-dev/git-ninja/app/git"
	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/permafrost-dev/git-ninja/app/utils"
	"github.com/spf13/cobra"
)

func getAllBranchDataSortedByAge(lines []string, availableBranches map[string]bool) []*git.BranchCheckoutInfo {
	var lineRegex = regexp.MustCompile(`([0-9]+) ~ (checkout):.+ ([^~]+) ~ HEAD@{(.*)}`)
	result := make([]*git.BranchCheckoutInfo, 0)

	for _, line := range lines {
		info := git.GetBranchInfoFromReflogLine(lineRegex, line, 4)

		if info != nil && utils.MapEntryExists(info.BranchName, availableBranches) && !git.SliceContainsBranchCommitData(result, info) {
			result = append(result, info)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.After(result[j].Timestamp)
	})

	return result
}

var flagRegex bool = false
var flagCheckoutFirst bool = false

var searchBranchesCmd = &cobra.Command{
	Use:   "branch:search [--regex|-r] <substring-or-regex>",
	Short: "Search branch names for matching substrings",
	Long:  `Searches branch names for matching substrings and displays a list of matching branches.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Error: No search string provided.")
			return
		}

		searchFor := args[0]
		lines, _ := git.GetGitReflogLines("%at ~ %gs ~ %gd")
		existingBranches, _ := helpers.GetAvailableBranchesMap()
		sortedBranches := getAllBranchDataSortedByAge(lines, existingBranches)

		var matches []*git.BranchCheckoutInfo

		for _, branchData := range sortedBranches {
			if !flagRegex && strings.Contains(branchData.BranchName, searchFor) {
				matches = append(matches, branchData)
			}

			if flagRegex && utils.StringMatchesRegexPattern(searchFor, branchData.BranchName) {
				matches = append(matches, branchData)
			}
		}

		if len(matches) == 0 {
			fmt.Println("No matching branches found.")
		}

		if len(matches) > 0 && flagCheckoutFirst {
			fmt.Printf("  \033[33m%-16s \033[37;1m %s\033[0m\n", matches[0].RelativeTime, matches[0].BranchName)
			helpers.RunCommandOnStdout("git", "checkout", matches[0].BranchName)
			return
		}

		for _, branch := range matches {
			fmt.Printf("  \033[33m%-16s \033[37;1m %s\033[0m\n", branch.RelativeTime, branch.BranchName)
		}
	},
}

func init() {
	rootCmd.AddCommand(searchBranchesCmd)

	searchBranchesCmd.Flags().BoolVarP(&flagRegex, "regex", "r", false, "Search using a regular expression pattern")
	searchBranchesCmd.Flags().BoolVarP(&flagCheckoutFirst, "checkout", "o", false, "Checkout the first matching branch")
}
