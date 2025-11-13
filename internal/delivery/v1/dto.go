package v1

import (
	"avito-internship/internal/domain"
)

// TODO: добавить валидацию
type TeamMember struct {
	Id       string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type User struct {
	Id       string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PullRequest struct {
	Id                string          `json:"pull_request_id"`
	Name              string          `json:"pull_request_name"`
	AuthorId          string          `json:"author_id"`
	Status            domain.PRStatus `json:"status"`
	AssignedReviewers []string        `json:"assigned_reviewers"`
	CreatedAt         string          `json:"createdAt"` // time.Time
	MergedAt          *string         `json:"mergedAt"`  // time.Time
}

type PullRequestShort struct {
	Id       string          `json:"pull_request_id"`
	Name     string          `json:"pull_request_name"`
	AuthorId string          `json:"author_id"`
	Status   domain.PRStatus `json:"status"`
}

type TeamAddRes struct {
	Team Team `json:"team"`
}

type SetIsActiveReq struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveRes struct {
	User User `json:"user"`
}

type CreatePullRequestReq struct {
	Id       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorId string `json:"author_id"`
}

type CreatePullRequestRes struct {
	PullRequest PullRequest `json:"pr"`
}

type PullRequestMergeReq struct {
	Id string `json:"pull_request_id"`
}

type PullRequestReassignReq struct {
	PullRequestId string `json:"pull_request_id"`
	OldUserId     string `json:"old_user_id"`
}

type PullRequestReassignRes struct {
	Pr         PullRequest `json:"pr"`
	ReplacedBy string      `json:"replaced_by"` // айди нового ревьюера
}

type GetReviewQueryReq struct {
	UserID string `form:"user_id" binding:"required"`
}

type GetReviewRes struct {
	UserId       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}
