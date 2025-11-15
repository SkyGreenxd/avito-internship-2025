package v1

import (
	"avito-internship/internal/domain"
)

type TeamMemberDTO struct {
	Id       string `json:"user_id" binding:"required,userid"`
	Username string `json:"username" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

type TeamDTO struct {
	TeamName string          `json:"team_name" binding:"required"`
	Members  []TeamMemberDTO `json:"members" binding:"required,dive"`
}

type UserDTO struct {
	Id       string `json:"user_id" binding:"required,userid"`
	Username string `json:"username" binding:"required"`
	TeamName string `json:"team_name" binding:"required"`
	IsActive bool   `json:"is_active" binding:"required"`
}

type PullRequestDTO struct {
	Id                string          `json:"pull_request_id" binding:"required,prid"`
	Name              string          `json:"pull_request_name" binding:"required"`
	AuthorId          string          `json:"author_id" binding:"required"`
	Status            domain.PRStatus `json:"status" binding:"required"`
	AssignedReviewers []string        `json:"assigned_reviewers" binding:"required"`
	CreatedAt         *string         `json:"createdAt" binding:"omitempty"`
	MergedAt          *string         `json:"mergedAt" binding:"omitempty"`
}

type PullRequestShort struct {
	Id       string          `json:"pull_request_id" binding:"required,prid"`
	Name     string          `json:"pull_request_name" binding:"required"`
	AuthorId string          `json:"author_id" binding:"required"`
	Status   domain.PRStatus `json:"status" binding:"required"`
}

type TeamAddReq struct {
	TeamName string          `json:"team_name" binding:"required"`
	Members  []TeamMemberDTO `json:"members" binding:"required,dive"`
}

type TeamAddRes struct {
	Team TeamDTO `json:"team"`
}

type GetTeamQueryReq struct {
	TeamName string `form:"team_name" binding:"required"`
}

type GetTeamRes struct {
	TeamName string          `json:"team_name"`
	Members  []TeamMemberDTO `json:"members"`
}

type SetIsActiveReq struct {
	UserId   string `json:"user_id" binding:"required,userid"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

type SetIsActiveRes struct {
	User UserDTO `json:"user"`
}

type CreatePullRequestReq struct {
	Id       string `json:"pull_request_id" binding:"required,prid"`
	Name     string `json:"pull_request_name" binding:"required"`
	AuthorId string `json:"author_id" binding:"required,userid"`
}

type CreatePullRequestDTO struct {
	Id                string          `json:"pull_request_id" binding:"required,prid"`
	Name              string          `json:"pull_request_name" binding:"required"`
	AuthorId          string          `json:"author_id" binding:"required"`
	Status            domain.PRStatus `json:"status" binding:"required"`
	AssignedReviewers []string        `json:"assigned_reviewers" binding:"required"`
}

type CreatePullRequestRes struct {
	PullRequest CreatePullRequestDTO `json:"pr"`
}

type PullRequestMergeReq struct {
	Id string `json:"pull_request_id" binding:"required,prid"`
}

type MergePullRequestDTO struct {
	Id                string          `json:"pull_request_id" binding:"required,prid"`
	Name              string          `json:"pull_request_name" binding:"required"`
	AuthorId          string          `json:"author_id" binding:"required"`
	Status            domain.PRStatus `json:"status" binding:"required"`
	AssignedReviewers []string        `json:"assigned_reviewers" binding:"required"`
	MergedAt          *string         `json:"mergedAt" binding:"omitempty"`
}

type PullRequestMergeRes struct {
	PullRequest MergePullRequestDTO `json:"pr"`
}

type PullRequestReassignReq struct {
	PullRequestId string `json:"pull_request_id" binding:"required,prid"`
	OldReviewerId string `json:"old_reviewer_id" binding:"required,userid"`
}

type PullRequestReassignRes struct {
	Pr         PullRequestDTO `json:"pr"`
	ReplacedBy string         `json:"replaced_by"`
}

type GetReviewQueryReq struct {
	UserID string `form:"user_id" binding:"required,userid"`
}

type GetReviewRes struct {
	UserId       string             `json:"user_id"`
	PullRequests []PullRequestShort `json:"pull_requests"`
}
