package e

import "fmt"

var (
	ErrTeamIsExists = fmt.Errorf("team with this name already exists")
	ErrTeamNotFound = fmt.Errorf("team not found")

	ErrUserNotFound = fmt.Errorf("user not found")

	PRIsExists = fmt.Errorf("pool request is exists")
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
