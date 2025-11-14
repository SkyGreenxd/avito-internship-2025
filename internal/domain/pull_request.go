package domain

import "time"

type PullRequest struct {
	Id                string
	Name              string
	AuthorId          string
	StatusId          int
	NeedMoreReviewers bool
	CreatedAt         time.Time
	MergedAt          *time.Time
}

func NewPoolRequest(id, name, authorId string, statusId int) *PullRequest {
	return &PullRequest{
		Id:                id,
		Name:              name,
		AuthorId:          authorId,
		StatusId:          statusId,
		NeedMoreReviewers: true,
		CreatedAt:         time.Now(),
	}
}
