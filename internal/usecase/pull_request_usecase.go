package usecase

import (
	"avito-internship/internal/domain"
	r "avito-internship/internal/repository"
	"avito-internship/pkg/e"
	"context"
	"slices"
)

const (
	maxReviewers = 2
)

type PullRequestUseCase struct {
	prRepo       r.PullRequestRepository
	reviewerRepo r.PrReviewerRepository
	teamRepo     r.TeamRepository
	userRepo     r.UserRepository
}

func NewPullRequestUseCase(prRepo r.PullRequestRepository, reviewerRepo r.PrReviewerRepository,
	teamRepo r.TeamRepository, userRepo r.UserRepository) *PullRequestUseCase {
	return &PullRequestUseCase{
		prRepo:       prRepo,
		reviewerRepo: reviewerRepo,
		teamRepo:     teamRepo,
		userRepo:     userRepo,
	}
}

// PullRequestCreate создает PR и назначает до maxReviewers ревьюеров
func (p *PullRequestUseCase) PullRequestCreate(ctx context.Context, req CreatePullRequestReq) (CreatePullRequestRes, error) {
	const op = "PullRequestUseCase.PullRequestCreate"

	pr := domain.NewPoolRequest(req.Id, req.Name, req.AuthorId)
	newPr, err := p.prRepo.Create(ctx, *pr)
	if err != nil {
		return CreatePullRequestRes{}, e.Wrap(op, err)
	}

	reviewers, err := p.userRepo.GetReviewCandidates(ctx, pr.AuthorId, maxReviewers)
	if err != nil {
		return CreatePullRequestRes{}, e.Wrap(op, err)
	}

	reviewersIds := make([]string, len(reviewers))
	for i := range reviewers {
		reviewersIds = append(reviewersIds, reviewers[i].Id)
	}

	if len(reviewers) > 0 {
		err := p.reviewerRepo.AddReviewers(ctx, newPr.Id, reviewersIds)
		if err != nil {
			return CreatePullRequestRes{}, e.Wrap(op, err)
		}
	}

	prDTO := NewPullRequestDTO(*pr, reviewersIds)
	return NewCreatePullRequestRes(prDTO), nil
}

// prRepo update
func (p *PullRequestUseCase) PullRequestMerge(ctx context.Context, req PullRequestMergeReq) (PullRequestMergeRes, error) {
	const op = "PullRequestUseCase.PullRequestMerge"

	updPr, reviewersIds, err := p.prRepo.SetMergedStatus(ctx, req.Id)
	if err != nil {
		return PullRequestMergeRes{}, e.Wrap(op, err)
	}

	prDTO := NewPullRequestDTO(updPr, reviewersIds)
	return NewPullRequestMergeRes(prDTO), nil
}

// reviewerRepo UpdateReviewer
func (p *PullRequestUseCase) ReviewerReassign(ctx context.Context, req PullRequestReassignReq) (PullRequestReassignRes, error) {
	const op = "PullRequestUseCase.PullRequestReassign"

	// пользователя нет
	_, err := p.userRepo.GetById(ctx, req.OldUserId)
	if err != nil {
		return PullRequestReassignRes{}, e.Wrap(op, err)
	}

	// pr нет
	pr, reviewersIds, err := p.prRepo.GetByPrIdWithReviewersIds(ctx, req.PullRequestId)
	if err != nil {
		return PullRequestReassignRes{}, e.Wrap(op, err)
	}

	oldReviewerIndex := slices.Index(reviewersIds, req.OldUserId)
	if oldReviewerIndex == -1 {
		return PullRequestReassignRes{}, e.Wrap(op, e.ErrPrReviewerNotAssigned)
	}

	// status == merged
	if pr.Status == domain.MERGED {
		return PullRequestReassignRes{}, e.Wrap(op, e.ErrPrMerged)
	}

	//нет доступных кандидатов
	candidates, err := p.userRepo.GetReassignCandidates(ctx, pr.AuthorId, req.OldUserId, maxReviewers)
	if err != nil {
		return PullRequestReassignRes{}, e.Wrap(op, err)
	}

	if len(candidates) == 0 {
		return PullRequestReassignRes{}, e.Wrap(op, e.ErrPrNoCandidate)
	}

	candidatesIds := make([]string, len(candidates))
	for i := range candidates {
		candidatesIds[i] = candidates[i].Id
	}

	newReviewerId, err := p.reviewerRepo.UpdateReviewer(ctx, req.OldUserId, candidatesIds[0], pr.Id)
	if err != nil {
		return PullRequestReassignRes{}, e.Wrap(op, err)
	}

	prDTO := NewPullRequestDTO(pr, reviewersIds)
	prDTO.AssignedReviewers[oldReviewerIndex] = newReviewerId

	return NewPullRequestReassignRes(prDTO, newReviewerId), nil
}
