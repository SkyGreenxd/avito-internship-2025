package postgres

import (
	"avito-internship/internal/domain"
	"avito-internship/pkg/e"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestsRepository struct {
	Pool *pgxpool.Pool
}

func NewPullRequestsRepository(pool *pgxpool.Pool) *PullRequestsRepository {
	return &PullRequestsRepository{Pool: pool}
}

func (p *PullRequestsRepository) Create(ctx context.Context, poolRequest *domain.PoolRequest) (*domain.PoolRequest, error) {
	const op = "PullRequestsRepository.Create"

	model := toPRModel(poolRequest)
	builder := sq.Insert("pull_requests").
		Columns("id", "name", "author_id", "status", "need_more_reviewers", "created_at").
		Values(model.Id, model.Name, model.AuthorId, model.Status, model.NeedMoreReviewers, model.CreatedAt).
		Suffix("RETURNING id, name, author_id, status, need_more_reviewers, created_at, merged_at")

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	err = p.Pool.QueryRow(ctx, query, args...).Scan(&model.Id, &model.Name, &model.AuthorId, &model.Status,
		&model.NeedMoreReviewers, &model.CreatedAt, &model.MergedAt)
	if err = postgresDuplicate(err, e.ErrPRIsExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainPR(model), nil
}

func (p *PullRequestsRepository) Update(ctx context.Context, poolRequest *domain.PoolRequest) (*domain.PoolRequest, error) {
	const op = "PullRequestsRepository.Update"

	model := toPRModel(poolRequest)
	builder := sq.Update("pull_requests").
		Set("status", model.Status).
		Where(sq.Eq{"id": model.Id}).
		Suffix("RETURNING id, name, author_id, status, need_more_reviewers, created_at, merged_at")

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	row := p.Pool.QueryRow(ctx, query, args...)
	if err := checkGetQueryResult(row.Scan(&model.Id, &model.Name, &model.AuthorId, &model.Status,
		&model.NeedMoreReviewers, &model.CreatedAt, &model.MergedAt), e.ErrPRNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainPR(model), nil
}

func toPRModel(p *domain.PoolRequest) *PoolRequestModel {
	return &PoolRequestModel{
		Id:                p.Id,
		Name:              p.Name,
		AuthorId:          p.AuthorId,
		Status:            p.Status,
		NeedMoreReviewers: p.NeedMoreReviewers,
		CreatedAt:         p.CreatedAt,
		MergedAt:          p.MergedAt,
	}
}

func toDomainPR(p *PoolRequestModel) *domain.PoolRequest {
	return &domain.PoolRequest{
		Id:                p.Id,
		Name:              p.Name,
		AuthorId:          p.AuthorId,
		Status:            p.Status,
		NeedMoreReviewers: p.NeedMoreReviewers,
		CreatedAt:         p.CreatedAt,
		MergedAt:          p.MergedAt,
	}
}
