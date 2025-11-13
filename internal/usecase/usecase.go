package usecase

import "context"

type UserUC interface {
	SetIsActive(ctx context.Context, req SetIsActiveReq) (SetIsActiveRes, error)
	GetReview(ctx context.Context, userId string) (GetReviewRes, error)
}

type TeamUC interface {
	AddTeam(ctx context.Context, req TeamDTO) (TeamAddRes, error)
	GetTeam(ctx context.Context, teamName string) (TeamDTO, error)
}

type PullRequestUC interface {
	PullRequestCreate(ctx context.Context, req CreatePullRequestReq) (CreatePullRequestRes, error)
	PullRequestMerge(ctx context.Context, req PullRequestMergeReq) (PullRequestMergeRes, error)
	ReviewerReassign(ctx context.Context, req PullRequestReassignReq) (PullRequestReassignRes, error)
}
