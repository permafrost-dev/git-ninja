package gitutils

import "time"

type BranchInfo struct {
	Name          string
	CheckoutCount int
	LastCheckout  time.Time
}
