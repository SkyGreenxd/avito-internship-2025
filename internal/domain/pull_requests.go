package domain

type PRStatus string

const (
	OPEN   PRStatus = "OPEN"
	MERGED PRStatus = "MERGED"
)

type PoolRequest struct {
	Id                int
	Name              string
	AuthorId          int
	Status            PRStatus
	NeedMoreReviewers bool
}

func NewPoolRequest(name string, authorId int) *PoolRequest {
	return &PoolRequest{
		Name:              name,
		AuthorId:          authorId,
		Status:            OPEN,
		NeedMoreReviewers: true,
	}
}
