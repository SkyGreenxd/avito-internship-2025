package pgdb

import (
	"avito-internship/internal/domain"
	"avito-internship/pkg/e"
	"avito-internship/pkg/transaction"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	Pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{Pool: pool}
}

func (u *UserRepository) UpdateIsActive(ctx context.Context, userId string, isActive bool) (domain.User, error) {
	const op = "UserRepository.UpdateIsActive"

	queryBuilder := sq.Update("users").
		Set("is_active", isActive).
		Where(sq.Eq{"id": userId}).
		Suffix("RETURNING id, name, is_active, team_id")

	query, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.User{}, e.Wrap(op, err)
	}

	var updModel UserModel
	err = u.Pool.QueryRow(ctx, query, args...).Scan(&updModel.Id, &updModel.Name, &updModel.IsActive, &updModel.TeamId)
	if err := checkGetQueryResult(err, e.ErrUserNotFound); err != nil {
		return domain.User{}, e.Wrap(op, err)
	}

	return toDomainUser(updModel), nil
}

func (u *UserRepository) GetById(ctx context.Context, userId string) (domain.User, error) {
	const op = "UserRepository.GetById"

	builder := sq.Select("id, name, is_active, team_id").
		From("users").
		Where(sq.Eq{"id": userId})

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.User{}, e.Wrap(op, err)
	}

	var model UserModel
	err = u.Pool.QueryRow(ctx, query, args...).Scan(&model.Id, &model.Name, &model.IsActive, &model.TeamId)
	if err := checkGetQueryResult(err, e.ErrUserNotFound); err != nil {
		return domain.User{}, e.Wrap(op, err)
	}

	return toDomainUser(model), nil
}

func (u *UserRepository) GetReviewCandidates(ctx context.Context, authorId string, maxReviewers int) ([]domain.User, error) {
	const op = "UserRepository.GetReviewCandidates"

	tx, err := transaction.TxFromCtx(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	query := `
       SELECT id, name, is_active, team_id
       FROM users
       WHERE
           team_id = (SELECT team_id FROM users WHERE id = $1)
           AND is_active = TRUE
           AND id != $1
       ORDER BY random()
       LIMIT $2
    `

	rows, err := tx.Query(ctx, query, authorId, maxReviewers)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer rows.Close()

	var models []UserModel
	for rows.Next() {
		var model UserModel
		if err := rows.Scan(&model.Id, &model.Name, &model.IsActive, &model.TeamId); err != nil {
			return nil, e.Wrap(op, err)
		}

		models = append(models, model)
	}

	if err := rows.Err(); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toArrDomainUser(models), nil
}

func (u *UserRepository) GetReassignCandidates(ctx context.Context, authorId string, excludeIds []string, maxCandidates int) ([]domain.User, error) {
	const op = "UserRepository.GetReassignCandidates"

	builder := sq.Select("id", "name", "is_active", "team_id").
		From("users").
		Where(sq.Expr("team_id = (SELECT team_id FROM users WHERE id = ?)", authorId)).
		Where(sq.Eq{"is_active": true}).
		Where(sq.NotEq{"id": excludeIds}).
		OrderBy("random()").
		Limit(uint64(maxCandidates))

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	rows, err := u.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer rows.Close()

	var candidates []UserModel
	for rows.Next() {
		var m UserModel
		if err := rows.Scan(&m.Id, &m.Name, &m.IsActive, &m.TeamId); err != nil {
			return nil, e.Wrap(op, err)
		}
		candidates = append(candidates, m)
	}

	if err := rows.Err(); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toArrDomainUser(candidates), nil
}

func (u *UserRepository) AddUsersToTeam(ctx context.Context, teamId int, users []domain.User) ([]domain.User, error) {
	const op = "UserRepository.AddUsersToTeam"

	tx, err := transaction.TxFromCtx(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	var userIDs []string
	var userNames []string
	var userActives []bool

	for _, user := range users {
		userIDs = append(userIDs, user.Id)
		userNames = append(userNames, user.Name)
		userActives = append(userActives, user.IsActive)
	}

	query := `
		INSERT INTO users (id, name, is_active, team_id)
		SELECT
			u.id,
			u.name,
			u.is_active,
			$4 -- team_id
		FROM
			UNNEST(
				$1::varchar[],
				$2::varchar[],
				$3::boolean[]
			) AS u(id, name, is_active)
		ON CONFLICT (id) DO UPDATE
		SET
			name = EXCLUDED.name,
			is_active = EXCLUDED.is_active,
			team_id = EXCLUDED.team_id
		RETURNING id, name, is_active, team_id
	`

	rows, err := tx.Query(ctx, query, userIDs, userNames, userActives, teamId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer rows.Close()

	updUsers := make([]domain.User, 0, len(users))
	for rows.Next() {
		var userId string
		var userName string
		var userIsActive bool
		var userTeamId int
		if err := rows.Scan(&userId, &userName, &userIsActive, &userTeamId); err != nil {
			return nil, e.Wrap(op, err)
		}

		updUser := domain.NewUser(userId, userName, userIsActive, userTeamId)
		updUsers = append(updUsers, *updUser)
	}

	if err := rows.Err(); err != nil {
		return nil, e.Wrap(op, err)
	}

	return updUsers, nil
}

func (u *UserRepository) DeactivateUsers(ctx context.Context, ids []string) ([]domain.User, error) {
	const op = "UserRepository.DeactivateTeamMembers"

	tx, err := transaction.TxFromCtx(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	queryBuilder := sq.Update("users").
		Set("is_active", false).
		Where(sq.Eq{"id": ids}).
		Suffix("RETURNING id, name, team_id, is_active")

	query, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer rows.Close()

	var users []UserModel
	for rows.Next() {
		var m UserModel
		if err := rows.Scan(&m.Id, &m.Name, &m.TeamId, &m.IsActive); err != nil {
			return nil, e.Wrap(op, err)
		}
		users = append(users, m)
	}

	if err := rows.Err(); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toArrDomainUser(users), nil
}

func toDomainUser(u UserModel) domain.User {
	return domain.User{
		Id:       u.Id,
		Name:     u.Name,
		IsActive: u.IsActive,
		TeamId:   u.TeamId,
	}
}

func toArrDomainUser(u []UserModel) []domain.User {
	users := make([]domain.User, 0, len(u))
	for _, user := range u {
		users = append(users, toDomainUser(user))
	}

	return users
}
