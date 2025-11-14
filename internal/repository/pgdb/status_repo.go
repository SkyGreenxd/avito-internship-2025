package pgdb

import (
	"avito-internship/internal/domain"
	"avito-internship/pkg/e"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatusRepo struct {
	Pool *pgxpool.Pool
}

func NewStatusRepo(pool *pgxpool.Pool) *StatusRepo {
	return &StatusRepo{
		Pool: pool,
	}
}

func (s *StatusRepo) GetById(ctx context.Context, statusId int) (domain.Status, error) {
	const op = "StatusRepo.GetById"

	query := `SELECT id, name FROM statuses WHERE id = $1`

	var model StatusModel
	err := s.Pool.QueryRow(ctx, query, statusId).Scan(&model.Id, &model.Name)
	if err = checkGetQueryResult(err, e.ErrStatusNotFound); err != nil {
		return domain.Status{}, err
	}

	return toDomainStatus(model), nil
}

func toDomainStatus(status StatusModel) domain.Status {
	return domain.Status{
		Id:   status.Id,
		Name: status.Name,
	}
}

func (s *StatusRepo) GetByName(ctx context.Context, statusName string) (domain.Status, error) {
	const op = "StatusRepo.GetByName"

	query := `SELECT id, name FROM statuses WHERE name = $1`

	var model StatusModel
	err := s.Pool.QueryRow(ctx, query, statusName).Scan(&model.Id, &model.Name)
	if err = checkGetQueryResult(err, e.ErrStatusNotFound); err != nil {
		return domain.Status{}, err
	}

	return toDomainStatus(model), nil
}
