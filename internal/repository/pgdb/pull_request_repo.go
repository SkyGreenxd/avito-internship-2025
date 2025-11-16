package pgdb

import (
	"avito-internship/internal/domain"
	r "avito-internship/internal/repository"
	"avito-internship/pkg/e"
	"avito-internship/pkg/transaction"
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

	tx, err := transaction.TxFromCtx(ctx)
	if err != nil {
		return domain.PullRequest{}, e.Wrap(op, err)
	}

	model := toPRModel(pullRequest)
	builder := sq.Insert("pull_requests").
		Columns("id", "name", "author_id", "status_id", "need_more_reviewers", "created_at").
		Values(model.Id, model.Name, model.AuthorId, model.StatusId, model.NeedMoreReviewers, model.CreatedAt).
		Suffix("RETURNING id, name, author_id, status_id, need_more_reviewers, created_at, merged_at")

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.PullRequest{}, e.Wrap(op, err)
	}

	err = tx.QueryRow(ctx, query, args...).Scan(&model.Id, &model.Name, &model.AuthorId, &model.StatusId, &model.NeedMoreReviewers, &model.CreatedAt, &model.MergedAt)
	err = postgresDuplicate(err, e.ErrPRIsExists)
	err = postgresForeignKeyViolation(err, e.ErrUserNotFound)
	if err != nil {
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

func (p *PullRequestsRepository) GetOpenPRsByReviewerIDs(ctx context.Context, reviewersIds []string, statusId int) (map[string]r.GetOpenPRsByReviewerIDsDTO, error) {
	const op = "PullRequestsRepository.GetOpenPRsByReviewerIDs"

	tx, err := transaction.TxFromCtx(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	query := `
       SELECT
          pr.id AS pr_id,
          pr.name AS pr_name,
          pr.author_id AS author_id,
          pr.status_id AS status_id,
          s.name AS status_name,
          pr.need_more_reviewers AS need_more_reviewers,
          pr.created_at AS created_at,
          pr.merged_at AS merged_at,
          r_all.reviewer_id AS reviewer_id
       FROM pull_requests pr
       JOIN pr_reviewers r_search ON pr.id = r_search.pr_id
       JOIN pr_reviewers r_all ON pr.id = r_all.pr_id
       JOIN statuses s ON pr.status_id = s.id
       WHERE r_search.reviewer_id = ANY($1)
          AND pr.status_id = $2
       ORDER BY pr.id;
    `

	rows, err := tx.Query(ctx, query, reviewersIds, statusId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer rows.Close()

	prTempMap := make(map[string]r.GetOpenPRsByReviewerIDsDTO)

	for rows.Next() {
		var prID, reviewerID, statusName string
		var pr domain.PullRequest

		if err := rows.Scan(
			&pr.Id,
			&pr.Name,
			&pr.AuthorId,
			&pr.StatusId,
			&statusName,
			&pr.NeedMoreReviewers,
			&pr.CreatedAt,
			&pr.MergedAt,
			&reviewerID,
		); err != nil {
			return nil, e.Wrap(op, err)
		}

		prID = pr.Id

		dto, exists := prTempMap[prID]
		if !exists {
			dto = r.GetOpenPRsByReviewerIDsDTO{
				Pr:           pr,
				ReviewersIds: []string{},
				StatusName:   statusName,
			}
			prTempMap[prID] = dto
		}

		found := false
		for _, id := range dto.ReviewersIds {
			if id == reviewerID {
				found = true
				break
			}
		}
		if !found {
			dto.ReviewersIds = append(dto.ReviewersIds, reviewerID)
			prTempMap[prID] = dto
		}
	}

	if err := rows.Err(); err != nil {
		return nil, e.Wrap(op, err)
	}

	return prTempMap, nil
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
