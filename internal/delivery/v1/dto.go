package v1

import (
	"avito-internship/internal/domain"
)

// TODO: добавить валидацию
type TeamMemberDTO struct {
	Id       string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type TeamDTO struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type UserDTO struct {
	Id       string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PullRequestDTO struct {
	Id                string          `json:"pull_request_id"`
	Name              string          `json:"pull_request_name"`
	AuthorId          string          `json:"author_id"`
	Status            domain.PRStatus `json:"status"`
	AssignedReviewers []string        `json:"assigned_reviewers"`
	CreatedAt         *string         `json:"createdAt"` // time.Time
	MergedAt          *string         `json:"mergedAt"`  // time.Time
}

type PullRequestShort struct {
	Id       string          `json:"pull_request_id"`
	Name     string          `json:"pull_request_name"`
	AuthorId string          `json:"author_id"`
	Status   domain.PRStatus `json:"status"`
}

type TeamAddReq struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type TeamAddRes struct {
	Team TeamDTO `json:"team"`
}

type GetTeamQueryReq struct {
	TeamName string `json:"team_name"`
}

type GetTeamRes struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type SetIsActiveReq struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveRes struct {
	User UserDTO `json:"user"`
}

type CreatePullRequestReq struct {
	Id       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorId string `json:"author_id"`
}

type CreatePullRequestRes struct {
	PullRequest PullRequestDTO `json:"pr"`
}

type PullRequestMergeReq struct {
	Id string `json:"pull_request_id"`
}

type PullRequestMergeRes struct {
	PullRequest PullRequestDTO `json:"pr"`
}

type PullRequestReassignReq struct {
	PullRequestId string `json:"pull_request_id"`
	OldUserId     string `json:"old_user_id"`
}

type PullRequestReassignRes struct {
	Pr         PullRequestDTO `json:"pr"`
	ReplacedBy string         `json:"replaced_by"` // айди нового ревьюера
}

type GetReviewQueryReq struct {
	UserID string `form:"user_id" binding:"required"`
}

type GetReviewRes struct {
	UserId       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}
