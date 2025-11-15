package e

import "fmt"

var (
	// Ошибки с командами
	ErrTeamIsExists = fmt.Errorf("team_name already exists")
	ErrTeamNotFound = fmt.Errorf("team not found")

	// Ошибки с пользователями
	ErrUserNotFound = fmt.Errorf("user not found")

	// Ошибки с PullRequest
	ErrPRIsExists            = fmt.Errorf("PR id already exists")
	ErrPRNotFound            = fmt.Errorf("pool request not found")
	ErrPrMerged              = fmt.Errorf("cannot reassign on merged PR")
	ErrPrReviewerNotAssigned = fmt.Errorf("reviewer is not assigned to this PR")
	ErrPrNoCandidate         = fmt.Errorf("no active replacement candidate in team")

	// Ошибки со статусом
	ErrStatusNotFound = fmt.Errorf("status not found")
	ErrInvalidStatus  = fmt.Errorf("invalid status")

	// Ошибки пользователя
	ErrInvalidRequestBody = fmt.Errorf("invalid request body")
	ErrResourceNotFound   = fmt.Errorf("resource not found")
	ErrUnauthorized       = fmt.Errorf("unauthorized")
	ErrEmptyMembers       = fmt.Errorf("member list is empty")

	// Ошибки сервера
	ErrInternalServerError = fmt.Errorf("internal server error")

	// Ошибки с валидацией
	ErrValidatorFailed = fmt.Errorf("validator failed")

	// Ошибки с транзакциями
	ErrTransactionNotFound = fmt.Errorf("transaction not found in context")
)

const (
	NOT_FOUND    = "NOT_FOUND"
	TEAM_EXISTS  = "TEAM_EXISTS"
	PR_EXISTS    = "PR_EXISTS"
	PR_MERGED    = "PR_MERGED"
	NOT_ASSIGNED = "NOT_ASSIGNED"
	NO_CANDIDATE = "NO_CANDIDATE"
	SERVER_ERR   = "SERVER_ERR"
	BAD_REQUEST  = "BAD_REQUEST"
)

func Wrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
