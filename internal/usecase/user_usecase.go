package usecase

import (
	r "avito-internship/internal/repository"
	"avito-internship/pkg/e"
	"context"
)

type UserUseCase struct {
	reviewerRepo r.PrReviewerRepository
	userRepo     r.UserRepository
	teamRepo     r.TeamRepository
}

func NewUserUseCase(reviewerRepo r.PrReviewerRepository, userRepo r.UserRepository, teamRepo r.TeamRepository) UserUseCase {
	return UserUseCase{
		reviewerRepo: reviewerRepo,
		userRepo:     userRepo,
		teamRepo:     teamRepo,
	}
}

// SetIsActive меняет статус активности пользователя
func (u *UserUseCase) SetIsActive(ctx context.Context, req SetIsActiveReq) (SetIsActiveRes, error) {
	const op = "UserUseCase.SetIsActive"

	updUser, err := u.userRepo.UpdateIsActive(ctx, req.UserId, req.IsActive)
	if err != nil {
		return SetIsActiveRes{}, e.Wrap(op, err)
	}

	team, err := u.teamRepo.GetTeamByUserId(ctx, updUser.Id)
	if err != nil {
		return SetIsActiveRes{}, e.Wrap(op, err)
	}

	return NewSetIsActiveRes(updUser.Id, updUser.Name, team.Name, updUser.IsActive), nil
}

// GetReview возвращает список PR'ов, где пользователь назначен ревьюером
func (u *UserUseCase) GetReview(ctx context.Context, userId string) (GetReviewRes, error) {
	const op = "UserUseCase.GetReview"

	prs, err := u.reviewerRepo.GetPRByReviewer(ctx, userId)
	if err != nil {
		return GetReviewRes{}, e.Wrap(op, err)
	}

	return NewGetReviewRes(userId, prs), nil
}
