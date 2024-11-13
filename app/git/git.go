package git

import (
	"regexp"
	"strings"
	"time"

	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/permafrost-dev/git-ninja/app/utils"
)

type BranchInfo struct {
	Name            string
	CheckoutCount   int
	CheckoutHistory []*BranchCheckoutInfo
	CheckedOutLast  time.Time
}

type BranchCheckoutInfo struct {
	BranchName   string
	RelativeTime string
	Timestamp    time.Time
}

func GetBranchInfoFromReflogLine(pattern *regexp.Regexp, reflogLine string, minMatchCount int) *BranchCheckoutInfo {
	matches := pattern.FindStringSubmatch(reflogLine)

	if matches == nil || len(matches) < minMatchCount {
		return nil
	}

	return &BranchCheckoutInfo{
		BranchName:   strings.TrimSpace(matches[3]),
		RelativeTime: strings.TrimSpace(matches[4]),
		Timestamp:    utils.ParseTimestampIntoTime(matches[1]),
	}
}

func GetGitReflogLines(prettyFmt string) ([]string, error) {
	output, err := helpers.RunCommandBuffered("git", "reflog", "show", "--pretty=format:'"+prettyFmt+"'", "--date=relative")
	if err != nil {
		return make([]string, 0), err
	}

	return strings.Split(output, "\n"), nil
}

func SliceContainsBranchCommitData(slice []*BranchCheckoutInfo, info *BranchCheckoutInfo) bool {
	for _, item := range slice {
		if item.BranchName == info.BranchName {
			return true
		}
	}

	return false
}
