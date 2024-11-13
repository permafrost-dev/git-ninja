package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/permafrost-dev/git-ninja/app/git"
	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/permafrost-dev/git-ninja/app/utils"
	"github.com/spf13/cobra"
)

var flagCount int = 10
var flagFilterIgnore string = ""
var lineRegex = regexp.MustCompile(`([0-9]+) ~ (checkout):.+ ([^~]+) ~ HEAD@{(.*)}`)
var refLogCheckoutsFmt = "%at ~ %gs ~ %gd"

var listRecentBranchesCmd = &cobra.Command{
	Use:   "branch:recent [--count|-c <count>]",
	Short: "Show recently checked out branch names",
	Run: func(c *cobra.Command, args []string) {
		existingBranches, _ := helpers.GetAvailableBranchesMap()
		currentBranch, _ := helpers.GetCurrentBranchName()
		lines, _ := git.GetGitReflogLines(refLogCheckoutsFmt)
		seen := make(map[string]bool)
		count := 0

		for _, line := range lines {
			info := git.GetBranchInfoFromReflogLine(lineRegex, line, 4)
			if info == nil {
				continue
			}

			// exclude branches that are not in the list of current branches, i.e. branches that have been deleted
			if !utils.MapEntryExists(info.BranchName, existingBranches) {
				continue
			}

			// is the branch name excluded by the exclude flag?
			if utils.StringMatchesRegexPattern(flagFilterIgnore, info.BranchName) {
				continue
			}

			// don't show the current branch
			if strings.EqualFold(info.BranchName, currentBranch) {
				continue
			}

			// have we already seen this branch?
			if seen[info.BranchName] {
				continue
			}

			seen[info.BranchName] = true
			fmt.Printf("  \033[33m%-15s \033[37;1m %s\033[0m\n", info.RelativeTime, info.BranchName)

			if count += 1; count >= flagCount {
				break
			}
		}

		// if no branches were found, show the current branch
		if count == 0 {
			fmt.Printf("  \033[33m%-15s \033[37;1m %s\033[0m\n", "now", currentBranch)
		}
	},
}

func init() {
	rootCmd.AddCommand(listRecentBranchesCmd)

	listRecentBranchesCmd.Flags().IntVarP(&flagCount, "count", "c", 10, "Limit the number of branches to display")
	listRecentBranchesCmd.Flags().StringVarP(&flagFilterIgnore, "exclude", "e", "", "Exclude branches that match the provided regex")
}
