package repository

import (
	"avito-internship/internal/domain"
	"context"
)

type UserRepository interface {
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	GetById(ctx context.Context, id string) (*domain.User, error)
	AddUsersToTeam(ctx context.Context, teamId string, users []*domain.User) ([]*domain.User, error)
}

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) (*domain.Team, error) // TODO создает только команду
	GetByTeamNameWithUsers(ctx context.Context, teamName string) (*domain.TeamWithUsers, error)
	GetByTeamNameWithActiveUsers(ctx context.Context, teamName string) (*domain.TeamWithUsers, error)
}

type PullRequestRepository interface {
	Create(ctx context.Context, poolRequest *domain.PoolRequest) (*domain.PoolRequest, error)
	Update(ctx context.Context, poolRequest *domain.PoolRequest) (*domain.PoolRequest, error)
}

type PrReviewerRepository interface {
	AddReviewers(ctx context.Context, poolRequestId string, reviewersId []string) error
	GetPRByReviewer(ctx context.Context, userId string) (*domain.PrReviewer, error)
	UpdateReviewer(ctx context.Context, oldUserId string, newUserId string, poolRequestId string) error
}
