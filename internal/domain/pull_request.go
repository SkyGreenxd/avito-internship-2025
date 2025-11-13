package domain

import "time"

type PRStatus string

const (
	OPEN   PRStatus = "OPEN"
	MERGED PRStatus = "MERGED"
)

type PullRequest struct {
	Id                string
	Name              string
	AuthorId          string
	Status            PRStatus
	NeedMoreReviewers bool
	CreatedAt         time.Time
	MergedAt          *time.Time
}

func NewPoolRequest(id, name, authorId string) *PullRequest {
	return &PullRequest{
		Id:                id,
		Name:              name,
		AuthorId:          authorId,
		Status:            OPEN,
		NeedMoreReviewers: true,
		CreatedAt:         time.Now(),
	}
}
