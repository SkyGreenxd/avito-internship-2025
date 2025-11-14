package v1

import (
	"avito-internship/internal/domain"
)

// TODO: добавить валидацию
type TeamMemberDTO struct {
	Id       string `json:"user_id" binding:"required,regexp=^u([1-9]|[1-9][0-9]|[1-9][0-9]{2})$"`
	Username string `json:"username" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}

type TeamDTO struct {
	TeamName string          `json:"team_name" binding:"required"`
	Members  []TeamMemberDTO `json:"members" binding:"required"`
}

type UserDTO struct {
	Id       string `json:"user_id" binding:"required,regexp=^u([1-9]|[1-9][0-9]|[1-9][0-9]{2})$"`
	Username string `json:"username" binding:"required"`
	TeamName string `json:"team_name" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}

type PullRequestDTO struct {
	Id                string          `json:"pull_request_id" binding:"required,regexp=^pr-(100[1-9]|10[1-9][0-9]|1[1-9][0-9]{2}|[2-9][0-9]{3})$"`
	Name              string          `json:"pull_request_name" binding:"required"`
	AuthorId          string          `json:"author_id" binding:"required"`
	Status            domain.PRStatus `json:"status" binding:"required"`
	AssignedReviewers []string        `json:"assigned_reviewers" binding:"required"`
	CreatedAt         *string         `json:"createdAt" binding:"omitempty"`
	MergedAt          *string         `json:"mergedAt" binding:"omitempty"`
}

type PullRequestShort struct {
	Id       string          `json:"pull_request_id" binding:"required,regexp=^pr-(100[1-9]|10[1-9][0-9]|1[1-9][0-9]{2}|[2-9][0-9]{3})$"`
	Name     string          `json:"pull_request_name" binding:"required"`
	AuthorId string          `json:"author_id" binding:"required"`
	Status   domain.PRStatus `json:"status" binding:"required"`
}

type TeamAddReq struct {
	TeamName string          `json:"team_name" binding:"required"`
	Members  []TeamMemberDTO `json:"members" binding:"required"`
}

type TeamAddRes struct {
	Team TeamDTO `json:"team"`
}

type GetTeamQueryReq struct {
	TeamName string `json:"team_name" binding:"required"`
}

type GetTeamRes struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type SetIsActiveReq struct {
	UserId   string `json:"user_id" binding:"required,regexp=^u([1-9]|[1-9][0-9]|[1-9][0-9]{2})$"`
	IsActive bool   `json:"is_active" binding:"required"`
}

type SetIsActiveRes struct {
	User UserDTO `json:"user"`
}

type CreatePullRequestReq struct {
	Id       string `json:"pull_request_id" binding:"required,regexp=^pr-(100[1-9]|10[1-9][0-9]|1[1-9][0-9]{2}|[2-9][0-9]{3})$"`
	Name     string `json:"pull_request_name" binding:"required"`
	AuthorId string `json:"author_id" binding:"required,regexp=^u([1-9]|[1-9][0-9]|[1-9][0-9]{2})$"`
}

type CreatePullRequestRes struct {
	PullRequest PullRequestDTO `json:"pr"`
}

type PullRequestMergeReq struct {
	Id string `json:"pull_request_id" binding:"required,regexp=^pr-(100[1-9]|10[1-9][0-9]|1[1-9][0-9]{2}|[2-9][0-9]{3})$"`
}

type PullRequestMergeRes struct {
	PullRequest PullRequestDTO `json:"pr"`
}

type PullRequestReassignReq struct {
	PullRequestId string `json:"pull_request_id" binding:"required,regexp=^pr-(100[1-9]|10[1-9][0-9]|1[1-9][0-9]{2}|[2-9][0-9]{3})$"`
	OldUserId     string `json:"old_user_id" binding:"required,regexp=^u([1-9]|[1-9][0-9]|[1-9][0-9]{2})$"`
}

type PullRequestReassignRes struct {
	Pr         PullRequestDTO `json:"pr"`
	ReplacedBy string         `json:"replaced_by"`
}

type GetReviewQueryReq struct {
	UserID string `form:"user_id" binding:"required,regexp=^u([1-9]|[1-9][0-9]|[1-9][0-9]{2})$"`
}

type GetReviewRes struct {
	UserId       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests" `
}
