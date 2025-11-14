package pgdb

import (
	"avito-internship/internal/domain"
	r "avito-internship/internal/repository"
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

func (p *PullRequestsRepository) Create(ctx context.Context, pullRequest domain.PullRequest) (domain.PullRequest, error) {
	const op = "PullRequestsRepository.Create"

	model := toPRModel(pullRequest)
	builder := sq.Insert("pull_requests").
		Columns("id", "name", "author_id", "status_id", "need_more_reviewers", "created_at").
		Values(model.Id, model.Name, model.AuthorId, model.StatusId, model.NeedMoreReviewers, model.CreatedAt).
		Suffix("RETURNING id, name, author_id, status_id, need_more_reviewers, created_at, merged_at")

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.PullRequest{}, e.Wrap(op, err)
	}

	err = p.Pool.QueryRow(ctx, query, args...).Scan(&model.Id, &model.Name, &model.AuthorId, &model.StatusId, &model.NeedMoreReviewers, &model.CreatedAt, &model.MergedAt)
	if err := postgresForeignKeyViolation(err, e.ErrUserNotFound); err != nil {
		return domain.PullRequest{}, e.Wrap(op, err)
	}

	if err = postgresDuplicate(err, e.ErrPRIsExists); err != nil {
		return domain.PullRequest{}, e.Wrap(op, err)
	}

	return toDomainPR(model), nil
}

func (p *PullRequestsRepository) SetMergedStatus(ctx context.Context, statusId int, prId string) (r.SetMergedStatusDTO, error) {
	const op = "PullRequestsRepository.SetMergedStatus"

	query := `
		WITH updated_pr AS (
			UPDATE pull_requests
			SET status_id = $1, merged_at = NOW()
			WHERE id = $2
			RETURNING id, name, author_id, status_id, need_more_reviewers, created_at, merged_at
		)
		SELECT u.id, u.name, u.author_id, u.status_id, u.need_more_reviewers, u.created_at, u.merged_at,
		       r.reviewer_id
		FROM updated_pr u
		LEFT JOIN pr_reviewers r ON r.pr_id = u.id
	`

	rows, err := p.Pool.Query(ctx, query, statusId, prId)
	if err != nil {
		return r.SetMergedStatusDTO{}, e.Wrap(op, err)
	}
	defer rows.Close()

	var upd PullRequestModel
	reviewers := make([]string, 0)

	for rows.Next() {
		var reviewerID *string
		if err := rows.Scan(&upd.Id, &upd.Name, &upd.AuthorId, &upd.StatusId,
			&upd.NeedMoreReviewers, &upd.CreatedAt, &upd.MergedAt, &reviewerID); err != nil {
			return r.SetMergedStatusDTO{}, e.Wrap(op, err)
		}
		if reviewerID != nil {
			reviewers = append(reviewers, *reviewerID)
		}
	}

	if err := rows.Err(); err != nil {
		return r.SetMergedStatusDTO{}, e.Wrap(op, err)
	}

	if upd.Id == "" {
		return r.SetMergedStatusDTO{}, e.Wrap(op, e.ErrPRNotFound)
	}

	return r.NewSetMergedStatusDTO(toDomainPR(upd), reviewers), nil
}

func (p *PullRequestsRepository) GetByPrIdWithReviewersIds(ctx context.Context, prId string) (r.GetByPrIdWithReviewersIdsDTO, error) {
	const op = "PullRequestsRepository.GetByPrIdWithReviewersIds"

	builder := sq.Select(
		"pr.id", "pr.name", "pr.author_id", "pr.status_id", "pr.need_more_reviewers", "pr.created_at", "pr.merged_at",
		"s.name AS status_name",
		"r.reviewer_id",
	).
		From("pull_requests AS pr").
		LeftJoin("pr_reviewers AS r ON r.pr_id = pr.id").
		LeftJoin("statuses AS s ON pr.status_id = s.id").
		Where(sq.Eq{"pr.id": prId})

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return r.GetByPrIdWithReviewersIdsDTO{}, e.Wrap(op, err)
	}

	rows, err := p.Pool.Query(ctx, query, args...)
	if err != nil {
		return r.GetByPrIdWithReviewersIdsDTO{}, e.Wrap(op, err)
	}
	defer rows.Close()

	var (
		model        PullRequestModel
		reviewersIds = make([]string, 0)
		statusName   domain.PRStatus
		prFound      = false
	)

	for rows.Next() {
		var reviewerId string
		if err := rows.Scan(
			&model.Id,
			&model.Name,
			&model.AuthorId,
			&model.StatusId,
			&model.NeedMoreReviewers,
			&model.CreatedAt,
			&model.MergedAt,
			&statusName,
			&reviewerId,
		); err != nil {
			return r.GetByPrIdWithReviewersIdsDTO{}, e.Wrap(op, err)
		}
		prFound = true

		if reviewerId != "" {
			reviewersIds = append(reviewersIds, reviewerId)
		}
	}

	if err := rows.Err(); err != nil {
		return r.GetByPrIdWithReviewersIdsDTO{}, e.Wrap(op, err)
	}

	if !prFound {
		return r.GetByPrIdWithReviewersIdsDTO{}, e.Wrap(op, e.ErrPRNotFound)
	}

	return r.NewGetByPrIdWithReviewersIdsDTO(toDomainPR(model), reviewersIds, statusName), nil
}

func toPRModel(p domain.PullRequest) PullRequestModel {
	return PullRequestModel{
		Id:                p.Id,
		Name:              p.Name,
		AuthorId:          p.AuthorId,
		StatusId:          p.StatusId,
		NeedMoreReviewers: p.NeedMoreReviewers,
		CreatedAt:         p.CreatedAt,
		MergedAt:          p.MergedAt,
	}
}

func toDomainPR(p PullRequestModel) domain.PullRequest {
	return domain.PullRequest{
		Id:                p.Id,
		Name:              p.Name,
		AuthorId:          p.AuthorId,
		StatusId:          p.StatusId,
		NeedMoreReviewers: p.NeedMoreReviewers,
		CreatedAt:         p.CreatedAt,
		MergedAt:          p.MergedAt,
	}
}
