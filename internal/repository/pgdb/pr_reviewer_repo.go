package pgdb

import (
	"avito-internship/internal/domain"
	r "avito-internship/internal/repository"
	"avito-internship/pkg/e"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PrReviewerRepository struct {
	Pool *pgxpool.Pool
}

func NewPrReviewerRepository(pool *pgxpool.Pool) *PrReviewerRepository {
	return &PrReviewerRepository{Pool: pool}
}

func (p *PrReviewerRepository) AddReviewers(ctx context.Context, poolRequestId string, reviewersId []string) error {
	const op = "PrReviewerRepository.AddReviewers"

	builder := sq.Insert("pr_reviewers").
		Columns("reviewer_id", "pr_id")

	for _, reviewerId := range reviewersId {
		builder = builder.Values(reviewerId, poolRequestId)
	}

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return e.Wrap(op, err)
	}

	_, err = p.Pool.Exec(ctx, query, args...)
	if err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

func (p *PrReviewerRepository) GetPRByReviewer(ctx context.Context, userId string) (r.GetPRByReviewerDTO, error) {
	const op = "PrReviewerRepository.GetByReviewer"

	builder := sq.Select(
		"pr.id",
		"pr.name",
		"pr.author_id",
		"pr.status_id",
		"pr.need_more_reviewers",
		"pr.created_at",
		"pr.merged_at",
		"s.name AS status_name",
	).
		From("pull_requests pr").
		Join("pr_reviewers r ON r.pr_id = pr.id").
		Join("statuses s ON s.id = pr.status_id").
		Where(sq.Eq{"r.reviewer_id": userId})

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return r.GetPRByReviewerDTO{}, e.Wrap(op, err)
	}

	rows, err := p.Pool.Query(ctx, query, args...)
	if err != nil {
		return r.GetPRByReviewerDTO{}, e.Wrap(op, err)
	}
	defer rows.Close()

	pullRequests := make([]domain.PullRequest, 0)
	statusNames := make([]domain.PRStatus, 0)

	for rows.Next() {
		var (
			pullRequest domain.PullRequest
			statusName  domain.PRStatus
		)
		if err := rows.Scan(
			&pullRequest.Id,
			&pullRequest.Name,
			&pullRequest.AuthorId,
			&pullRequest.StatusId,
			&pullRequest.NeedMoreReviewers,
			&pullRequest.CreatedAt,
			&pullRequest.MergedAt,
			&statusName,
		); err != nil {
			return r.GetPRByReviewerDTO{}, e.Wrap(op, err)
		}

		pullRequests = append(pullRequests, pullRequest)
		statusNames = append(statusNames, statusName)
	}

	if err := rows.Err(); err != nil {
		return r.GetPRByReviewerDTO{}, e.Wrap(op, err)
	}

	return r.NewGetPRByReviewerDTO(pullRequests, statusNames), nil
}

func (p *PrReviewerRepository) UpdateReviewer(ctx context.Context, oldUserId string, newUserId string, poolRequestId string) (string, error) {
	const op = "PrReviewerRepository.UpdateReviewer"

	builder := sq.Update("pr_reviewers").
		Set("reviewer_id", newUserId).
		Where(sq.Eq{
			"reviewer_id": oldUserId,
			"pr_id":       poolRequestId,
		}).
		Suffix("RETURNING reviewer_id")

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return "", e.Wrap(op, err)
	}

	var returnedPrID string
	err = p.Pool.QueryRow(ctx, query, args...).Scan(&returnedPrID)
	if err != nil {
		return "", e.Wrap(op, err)
	}

	return returnedPrID, nil
}
