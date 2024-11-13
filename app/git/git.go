package git

import (
	"regexp"
	"strings"
	"time"

	"github.com/permafrost-dev/git-ninja/app/helpers"
	"github.com/permafrost-dev/git-ninja/app/utils"
)

type BranchInfo struct {
	Name                 string
	CheckoutRelativeTime string
	CheckoutTimestamp    time.Time
}

func GetBranchInfoFromReflogLine(pattern *regexp.Regexp, reflogLine string) *BranchInfo {
	matches := pattern.FindStringSubmatch(reflogLine)

	if matches == nil || len(matches) < 4 {
		return nil
	}

	return &BranchInfo{
		Name:                 strings.TrimSpace(matches[3]),
		CheckoutRelativeTime: strings.TrimSpace(matches[4]),
		CheckoutTimestamp:    utils.ParseTimestampIntoTime(matches[1]),
	}
}

func GetGitReflogLines(prettyFmt string) ([]string, error) {
	output, err := helpers.RunCommandBuffered("git", "reflog", "show", "--pretty=format:'"+prettyFmt+"'", "--date=relative")
	if err != nil {
		return make([]string, 0), err
	}

	return strings.Split(output, "\n"), nil
}
