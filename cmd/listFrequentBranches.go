package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/permafrost-dev/git-ninja/app/gitutils"
	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/permafrost-dev/git-ninja/app/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "branch:freq",
		Short: "Show recently checked out branch names",
		Run: func(c *cobra.Command, args []string) {

			availableBranches, err := helpers.GetAvailableBranchesMap()
			if err != nil {
				fmt.Println("Error fetching branches:", err)
				return
			}

			cmd := exec.Command("git", "reflog", "show", "--pretty=format:%gs~%ci")
			var out bytes.Buffer
			cmd.Stdout = &out
			err = cmd.Run()
			if err != nil {
				fmt.Println("Error executing git command:", err)
				return
			}

			// Convert command output to string and split by new lines
			output := out.String()
			lines := strings.Split(output, "\n")

			// Prepare regular expressions to match 'checkout:' and extract the branch name and ISO date
			checkoutRegexp := regexp.MustCompile(`^checkout: moving from [^ ]+ to ([^ ]+)~(.*)$`)

			branchData := make(map[string]gitutils.BranchInfo)

			// Time threshold for very recent checkouts (last 2 days)
			recentThreshold := time.Now().AddDate(0, 0, -3)

			// Iterate through each line and apply the regex filters
			for _, line := range lines {
				if matches := checkoutRegexp.FindStringSubmatch(line); len(matches) == 3 {
					branch := matches[1]
					isoDate := strings.TrimSpace(matches[2])

					checkoutDate, err := time.Parse("2006-01-02 15:04:05 -0700", isoDate)
					if err != nil {
						fmt.Println("Error parsing date:", err)
						continue
					}

					if utils.MapEntryExists(branch, availableBranches) {
						info := branchData[branch]
						info.CheckoutCount++

						// Update latest checkout date
						if info.LastCheckout.IsZero() || checkoutDate.After(info.LastCheckout) {
							info.LastCheckout = checkoutDate
						}

						branchData[branch] = info
					}
				}
			}

			// Separate branches into very recent and older groups
			var veryRecentBranches, otherBranches []gitutils.BranchInfo
			for branch, info := range branchData {
				if info.LastCheckout.After(recentThreshold) {
					veryRecentBranches = append(veryRecentBranches, gitutils.BranchInfo{Name: branch, CheckoutCount: info.CheckoutCount, LastCheckout: info.LastCheckout})
				} else {
					otherBranches = append(otherBranches, gitutils.BranchInfo{Name: branch, CheckoutCount: info.CheckoutCount, LastCheckout: info.LastCheckout})
				}
			}

			// Sort very recent branches by recency (latest first)
			sort.Slice(veryRecentBranches, func(i, j int) bool {
				return veryRecentBranches[i].LastCheckout.After(veryRecentBranches[j].LastCheckout)
			})

			// Sort other branches by frequency and then by recency
			sort.Slice(otherBranches, func(i, j int) bool {
				if otherBranches[i].CheckoutCount == otherBranches[j].CheckoutCount {
					return otherBranches[i].LastCheckout.After(otherBranches[j].LastCheckout)
				}
				return otherBranches[i].CheckoutCount > otherBranches[j].CheckoutCount
			})

			// Combine the very recent and other branches, prioritizing very recent ones
			displayedBranches := append(veryRecentBranches, otherBranches...)

			// Display the sorted branches, ensuring at least 20 are shown
			count := 0
			for _, branch := range displayedBranches {
				infoStr := fmt.Sprintf("%2d checkouts, %-15s", branch.CheckoutCount, utils.GetRelativeTime(branch.LastCheckout))
				fmt.Printf("  \033[33m%28s \033[37;1m %s\033[0m\n", infoStr, branch.Name)
				count++
				if count >= 20 {
					break
				}
			}
		},
	})
}
