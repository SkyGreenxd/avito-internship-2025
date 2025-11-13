package postgres

import (
	"avito-internship/internal/domain"
	"time"
)

type UserModel struct {
	Id       string `db:"id"`
	Name     string `db:"name"`
	IsActive bool   `db:"is_active"`
	TeamId   *int   `db:"team_id"`
}

type TeamModel struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

type PullRequestModel struct {
	Id                string          `db:"id"`
	Name              string          `db:"name"`
	AuthorId          string          `db:"author_id"`
	Status            domain.PRStatus `db:"status"`
	NeedMoreReviewers bool            `db:"need_more_reviewers"`
	CreatedAt         time.Time       `db:"created_at"`
	MergedAt          *time.Time      `db:"merged_at"`
}
