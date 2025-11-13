package e

import "fmt"

var (
	ErrTeamIsExists = fmt.Errorf("team with this name already exists")
	ErrTeamNotFound = fmt.Errorf("team not found")

	ErrUserNotFound = fmt.Errorf("user not found")

	PRIsExists               = fmt.Errorf("pool request is exists")
	ErrPRNotFound            = fmt.Errorf("pool request not found")
	ErrPrMerged              = fmt.Errorf("cannot reassign on merged PR")
	ErrPrReviewerNotAssigned = fmt.Errorf("reviewer is not assigned to this PR")
	ErrPrNoCandidate         = fmt.Errorf("no active replacement candidate in team")
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
