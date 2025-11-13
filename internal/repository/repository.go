package repository

import (
	"avito-internship/internal/domain"
	"context"
)

type UserRepository interface {
	UpdateIsActive(ctx context.Context, userId string, isActive bool) (domain.User, error)
	GetById(ctx context.Context, userId string) (domain.User, error)
	GetReviewCandidates(ctx context.Context, authorId string, maxCandidates int) ([]domain.User, error)
	GetReassignCandidates(ctx context.Context, authorId, oldReviewerId string, maxCandidates int) ([]domain.User, error)
}

type TeamRepository interface {
	Create(ctx context.Context, team domain.Team) (domain.Team, error) // TODO создает только команду
	GetMembersByTeamNameWithUsers(ctx context.Context, teamName string) ([]domain.User, error)
	GetTeamByUserId(ctx context.Context, userId string) (domain.Team, error)
	AddUsersToTeam(ctx context.Context, teamId int, users []domain.User) ([]domain.User, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, pullRequest domain.PullRequest) (domain.PullRequest, error)
	SetMergedStatus(ctx context.Context, prId string) (domain.PullRequest, []string, error)
	GetByPrIdWithReviewersIds(ctx context.Context, prId string) (domain.PullRequest, []string, error)
}

type PrReviewerRepository interface {
	AddReviewers(ctx context.Context, pullRequestId string, reviewersId []string) error
	GetPRByReviewer(ctx context.Context, userId string) ([]domain.PullRequest, error)
	UpdateReviewer(ctx context.Context, oldUserId string, newUserId string, pullRequestId string) (string, error)
}
