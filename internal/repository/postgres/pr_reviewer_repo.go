package postgres

import (
	"avito-internship/internal/domain"
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

func (p *PrReviewerRepository) GetPRByReviewer(ctx context.Context, userId string) (*domain.PrReviewer, error) {
	const op = "PrReviewerRepository.GetByReviewer"

	builder := sq.Select("id", "name", "author_id", "status", "need_more_reviewers", "created_at", "merged_at").
		From("pull_requests pr").
		Join("pr_reviewers r ON r.pr_id = pr.id").
		Where(sq.Eq{"reviewer_id": userId})

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	rows, err := p.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer rows.Close()

	poolRequests := make([]*domain.PoolRequest, 0)
	for rows.Next() {
		var poolRequest domain.PoolRequest
		if err := rows.Scan(&poolRequest.Id, &poolRequest.Name, &poolRequest.AuthorId, &poolRequest.Status, &poolRequest.NeedMoreReviewers, &poolRequest.CreatedAt, &poolRequest.MergedAt); err != nil {
			return nil, e.Wrap(op, err)
		}

		poolRequests = append(poolRequests, &poolRequest)
	}

	if err := rows.Err(); err != nil {
		return nil, e.Wrap(op, err)
	}

	return domain.NewPrReviewer(userId, poolRequests), nil
}

func (p *PrReviewerRepository) UpdateReviewer(ctx context.Context, oldUserId string, newUserId string, poolRequestId string) error {
	const op = "PrReviewerRepository.UpdateReviewer"

	builder := sq.Update("pr_reviewers").
		Set("reviewer_id", newUserId).
		Where(sq.Eq{
			"reviewer_id": oldUserId,
			"pr_id":       poolRequestId,
		})

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
