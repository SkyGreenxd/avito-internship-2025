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
	userRepo     r.UserRepository
	statusRepo   r.StatusRepository
}

func NewPullRequestUseCase(prRepo r.PullRequestRepository, reviewerRepo r.PrReviewerRepository,
	userRepo r.UserRepository, statusRepo r.StatusRepository) *PullRequestUseCase {
	return &PullRequestUseCase{
		prRepo:       prRepo,
		reviewerRepo: reviewerRepo,
		userRepo:     userRepo,
		statusRepo:   statusRepo,
	}
}

// PullRequestCreate создает PR и назначает до maxReviewers ревьюеров
func (p *PullRequestUseCase) PullRequestCreate(ctx context.Context, req CreatePullRequestReq) (CreatePullRequestRes, error) {
	const op = "PullRequestUseCase.PullRequestCreate"

	status, err := p.statusRepo.GetByName(ctx, string(domain.OPEN))
	if err != nil {
		return CreatePullRequestRes{}, e.Wrap(op, err)
	}

	pr := domain.NewPoolRequest(req.Id, req.Name, req.AuthorId, status.Id)
	newPr, err := p.prRepo.Create(ctx, *pr)
	if err != nil {
		return CreatePullRequestRes{}, e.Wrap(op, err)
	}

	reviewers, err := p.userRepo.GetReviewCandidates(ctx, pr.AuthorId, maxReviewers)
	if err != nil {
		return CreatePullRequestRes{}, e.Wrap(op, err)
	}

	reviewersIds := make([]string, 0, len(reviewers))
	for i := range reviewers {
		reviewersIds = append(reviewersIds, reviewers[i].Id)
	}

	if len(reviewers) > 0 {
		err := p.reviewerRepo.AddReviewers(ctx, newPr.Id, reviewersIds)
		if err != nil {
			return CreatePullRequestRes{}, e.Wrap(op, err)
		}
	}

	prDTO := NewPullRequestDTO(*pr, reviewersIds, status.Name)
	return NewCreatePullRequestRes(prDTO), nil
}

func (p *PullRequestUseCase) PullRequestMerge(ctx context.Context, req PullRequestMergeReq) (PullRequestMergeRes, error) {
	const op = "PullRequestUseCase.PullRequestMerge"

	status, err := p.statusRepo.GetByName(ctx, string(domain.MERGED))
	if err != nil {
		return PullRequestMergeRes{}, e.Wrap(op, err)
	}

	dto, err := p.prRepo.SetMergedStatus(ctx, status.Id, req.Id)
	if err != nil {
		return PullRequestMergeRes{}, e.Wrap(op, err)
	}

	prDTO := NewPullRequestDTO(dto.Pr, dto.ReviewersIds, status.Name)
	return NewPullRequestMergeRes(prDTO), nil
}

func (p *PullRequestUseCase) ReviewerReassign(ctx context.Context, req PullRequestReassignReq) (PullRequestReassignRes, error) {
	const op = "PullRequestUseCase.PullRequestReassign"

	_, err := p.userRepo.GetById(ctx, req.OldReviewerId)
	if err != nil {
		return PullRequestReassignRes{}, e.Wrap(op, err)
	}

	dto, err := p.prRepo.GetByPrIdWithReviewersIds(ctx, req.PullRequestId)
	if err != nil {
		return PullRequestReassignRes{}, e.Wrap(op, err)
	}

	oldReviewerIndex := slices.Index(dto.ReviewersIds, req.OldReviewerId)
	if oldReviewerIndex == -1 {
		return PullRequestReassignRes{}, e.Wrap(op, e.ErrPrReviewerNotAssigned)
	}

	if dto.StatusName == domain.MERGED {
		return PullRequestReassignRes{}, e.Wrap(op, e.ErrPrMerged)
	}

	excludeIds := dto.ReviewersIds
	excludeIds = append(excludeIds, dto.Pr.AuthorId)
	candidates, err := p.userRepo.GetReassignCandidates(ctx, dto.Pr.AuthorId, excludeIds, maxReviewers)
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

	newReviewerId, err := p.reviewerRepo.UpdateReviewer(ctx, req.OldReviewerId, candidatesIds[0], dto.Pr.Id)
	if err != nil {
		return PullRequestReassignRes{}, e.Wrap(op, err)
	}

	prDTO := NewPullRequestDTO(dto.Pr, dto.ReviewersIds, dto.StatusName)
	prDTO.AssignedReviewers[oldReviewerIndex] = newReviewerId

	return NewPullRequestReassignRes(prDTO, newReviewerId), nil
}
