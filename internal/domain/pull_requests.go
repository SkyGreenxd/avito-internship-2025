package domain

import "time"

type PRStatus string

const (
	OPEN   PRStatus = "OPEN"
	MERGED PRStatus = "MERGED"
)

type PoolRequest struct {
	Id                string
	Name              string
	AuthorId          string
	Status            PRStatus
	NeedMoreReviewers bool
	CreatedAt         time.Time
	MergedAt          *time.Time
}

type PrReviewer struct {
	ReviewerId   string
	PoolRequests []*PoolRequest
}

func NewPoolRequest(id, name, authorId string) *PoolRequest {
	return &PoolRequest{
		Id:                id,
		Name:              name,
		AuthorId:          authorId,
		Status:            OPEN,
		NeedMoreReviewers: true,
		CreatedAt:         time.Now(),
	}
}

func NewPrReviewer(id string, prs []*PoolRequest) *PrReviewer {
	return &PrReviewer{
		ReviewerId:   id,
		PoolRequests: prs,
	}
}
