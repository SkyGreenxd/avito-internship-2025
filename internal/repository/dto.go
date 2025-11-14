package repository

import "avito-internship/internal/domain"

type SetMergedStatusDTO struct {
	Pr           domain.PullRequest
	ReviewersIds []string
}

type GetByPrIdWithReviewersIdsDTO struct {
	Pr           domain.PullRequest
	ReviewersIds []string
	StatusName   domain.PRStatus
}

type PrWithStatusName struct {
	Pr         domain.PullRequest
	StatusName domain.PRStatus
}

type GetPRByReviewerDTO struct {
	Prs []PrWithStatusName
}

func NewPrWithStatusName(pr domain.PullRequest, statusName domain.PRStatus) PrWithStatusName {
	return PrWithStatusName{
		Pr:         pr,
		StatusName: statusName,
	}
}

func NewArrPrWithStatusName(prs []domain.PullRequest, statusNames []domain.PRStatus) []PrWithStatusName {
	result := make([]PrWithStatusName, 0, len(prs))
	for idx, pr := range prs {
		result = append(result, NewPrWithStatusName(pr, statusNames[idx]))
	}

	return result
}

func NewGetPRByReviewerDTO(prs []domain.PullRequest, statusNames []domain.PRStatus) GetPRByReviewerDTO {
	return GetPRByReviewerDTO{
		Prs: NewArrPrWithStatusName(prs, statusNames),
	}
}

func NewSetMergedStatusDTO(pr domain.PullRequest, reviewersIds []string) SetMergedStatusDTO {
	return SetMergedStatusDTO{
		Pr:           pr,
		ReviewersIds: reviewersIds,
	}
}

func NewGetByPrIdWithReviewersIdsDTO(pr domain.PullRequest, reviewersIds []string, statusName domain.PRStatus) GetByPrIdWithReviewersIdsDTO {
	return GetByPrIdWithReviewersIdsDTO{
		Pr:           pr,
		ReviewersIds: reviewersIds,
		StatusName:   statusName,
	}
}
