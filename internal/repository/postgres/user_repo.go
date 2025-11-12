package postgres

import (
	"avito-internship/internal/domain"
	"avito-internship/pkg/e"
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

// Update обновляет is_active и team_id пользователя.
func (u *UserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	const op = "UserRepository.Update"

	model := toUserModel(user)

	queryBuilder := sq.Update("users").
		Set("is_active", model.IsActive).
		Set("team_id", model.TeamId).
		Where(sq.Eq{"id": model.Id}).
		Suffix("RETURNING id, name, is_active, team_id")

	query, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	row := u.Pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&model.Id, &model.Name, &model.IsActive, &model.TeamId); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainUser(model), nil
}

func (u *UserRepository) GetById(ctx context.Context, id string) (*domain.User, error) {
	const op = "UserRepository.GetById"

	builder := sq.Select("id, name, is_active, team_id").
		From("users").
		Where(sq.Eq{"id": id})

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	var model UserModel
	err = u.Pool.QueryRow(ctx, query, args...).Scan(&model.Id, &model.Name, &model.IsActive, &model.TeamId)
	if err := checkGetQueryResult(err, e.ErrUserNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainUser(&model), nil
}

func (u *UserRepository) AddUsersToTeam(ctx context.Context, teamId string, users []*domain.User) ([]*domain.User, error) {
	const op = "UserRepository.AddUsersToTeam"

	var userIDs []string
	var userNames []string
	var userActives []bool

	for _, user := range users {
		// TODO: добавить защиту что все три массива одной величины в USECASE
		userIDs = append(userIDs, user.Id)
		userNames = append(userNames, user.Name)
		userActives = append(userActives, true)
	}

	tx, err := u.Pool.Begin(ctx)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer tx.Rollback(ctx) // TODO: мб стоит обработать ошибку

	query := `
			INSERT INTO users (id, name, is_active, team_id)
			SELECT
				u.id,
				u.name,
				u.is_active,
				$4 -- team_id (один для всех)
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
			RETURNING id, name, is_active, team_id;
	`

	rows, err := tx.Query(ctx, query, userIDs, userNames, userActives, teamId)
	if err != nil {
		// TODO: какие то мб еще ошибки
		return nil, e.Wrap(op, err)
	}
	defer rows.Close()

	updUsers := make([]*domain.User, 0, len(users))
	for rows.Next() {
		var userId string
		var userName string
		var userIsActive bool
		var userTeamId *string
		if err := rows.Scan(&userId, &userName, &userIsActive, &userTeamId); err != nil {
			return nil, e.Wrap(op, err)
		}

		updUser := domain.NewUser(userId, userName, userIsActive, userTeamId)
		updUsers = append(updUsers, updUser)
	}

	if err := rows.Err(); err != nil {
		return nil, e.Wrap(op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, e.Wrap(op, err)
	}

	return updUsers, nil
}

func toUserModel(u *domain.User) *UserModel {
	return &UserModel{
		Id:       u.Id,
		Name:     u.Name,
		IsActive: u.IsActive,
		TeamId:   u.TeamId,
	}
}

func toDomainUser(u *UserModel) *domain.User {
	return &domain.User{
		Id:       u.Id,
		Name:     u.Name,
		IsActive: u.IsActive,
		TeamId:   u.TeamId,
	}
}
