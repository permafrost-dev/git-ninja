package cmd

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/permafrost-dev/git-ninja/app/git"
	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/permafrost-dev/git-ninja/app/utils"
	"github.com/permafrost-dev/git-ninja/lib/integrations/jira"
	"github.com/spf13/cobra"
)

var flagCount int = 10
var flagJira bool = false
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
		seen := make(map[string]*git.BranchCheckoutInfo)
		count := 0

		jiraIssues := make([]string, 0)

		if flagJira {
			if len(os.Getenv("JIRA_SUBDOMAIN")) == 0 || len(os.Getenv("JIRA_EMAIL_ADDRESS")) == 0 || len(os.Getenv("JIRA_API_TOKEN")) == 0 {
				fmt.Println("Error: JIRA_SUBDOMAIN, JIRA_EMAIL_ADDRESS and JIRA_API_TOKEN environment variables must be set.")
				return
			}

			jiraIssues = jira.GetJiraTicketIDs(os.Getenv("JIRA_SUBDOMAIN"), os.Getenv("JIRA_EMAIL_ADDRESS"))
		}

		sorted := make([]git.BranchInfo, 0)

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
			if seen[info.BranchName] != nil {
				continue
			}

			var rank int64 = info.Timestamp.Unix() - 1610000000

			if flagJira {
				for idx, issue := range jiraIssues {
					jiraHash, _ := jira.HashJiraIssueKey(info.BranchName)

					if strings.Contains(info.BranchName, issue) {
						rank += int64(1000*(len(jiraIssues)-idx)) + (jiraHash * -50)
					} else {
						rank -= ((1000 + jiraHash) + int64(100*(len(jiraIssues)-idx))) * 4
					}
				}
			}

			rank = rank / 100

			seen[info.BranchName] = info
			sorted = append(sorted, git.BranchInfo{Rank: rank, Name: info.BranchName, CheckedOutLast: info.Timestamp, CheckoutHistory: []*git.BranchCheckoutInfo{info}})
		}

		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Rank >= sorted[j].Rank
		})

		for _, bi := range sorted {
			fmt.Printf("  \033[33m%-15s %-5d \033[37;1m %s\033[0m\n", utils.GetRelativeTime(bi.CheckedOutLast), bi.Rank, bi.Name)

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
	listRecentBranchesCmd.Flags().BoolVarP(&flagJira, "jira", "J", false, "Use JIRA issues to help rank branches")
}
