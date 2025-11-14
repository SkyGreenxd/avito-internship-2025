package pgdb

import (
	"avito-internship/internal/domain"
	"avito-internship/pkg/e"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	Pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{Pool: pool}
}

func (t *TeamRepository) Create(ctx context.Context, team domain.Team) (domain.Team, error) {
	const op = "TeamRepository.Create"

	model := toTeamModel(team)
	queryBuilder := sq.Insert("teams").
		Columns("name").
		Values(model.Name).
		Suffix("RETURNING id, name")

	query, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.Team{}, e.Wrap(op, err)
	}

	err = t.Pool.QueryRow(ctx, query, args...).Scan(&model.Id, &model.Name)
	if err = postgresDuplicate(err, e.ErrTeamIsExists); err != nil {
		return domain.Team{}, e.Wrap(op, err)
	}

	return toDomainTeam(model), nil
}

func (t *TeamRepository) GetMembersByTeamNameWithUsers(ctx context.Context, teamName string) ([]domain.User, error) {
	const op = "TeamRepository.GetByTeamNameWithUsers"

	builder := sq.Select(
		"users.id", "users.name", "users.is_active", "users.team_id",
	).
		From("teams").
		LeftJoin("users ON teams.id = users.team_id").
		Where(sq.Eq{"teams.name": teamName})

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	rows, err := t.Pool.Query(ctx, query, args...)
	if err := checkGetQueryResult(err, e.ErrTeamNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}
	defer rows.Close()

	var (
		usersModel = make([]UserModel, 0)
		teamFound  bool
	)

	for rows.Next() {
		teamFound = true

		var (
			uId       *string
			uName     *string
			uIsActive *bool
			uTeamId   *int
		)

		err := rows.Scan(
			&uId, &uName, &uIsActive, &uTeamId,
		)
		if err != nil {
			return nil, e.Wrap(op, err)
		}

		if uId != nil {
			usersModel = append(usersModel, UserModel{
				Id:       *uId,
				Name:     *uName,
				IsActive: *uIsActive,
				TeamId:   uTeamId,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, e.Wrap(op, err)
	}

	if !teamFound {
		return nil, e.Wrap(op, e.ErrTeamNotFound)
	}

	return toArrDomainUser(usersModel), nil
}

func (t *TeamRepository) GetTeamByUserId(ctx context.Context, userId string) (domain.Team, error) {
	const op = "TeamRepository.GetTeamByUserId"

	builder := sq.Select("teams.id", "teams.name").
		From("teams").
		Join("users ON teams.id = users.team_id").
		Where(sq.Eq{"users.id": userId}).
		Limit(1)

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return domain.Team{}, e.Wrap(op, err)
	}

	var model TeamModel
	err = t.Pool.QueryRow(ctx, query, args...).Scan(&model.Id, &model.Name)
	if err := checkGetQueryResult(err, e.ErrUserNotFound); err != nil {
		return domain.Team{}, e.Wrap(op, err)
	}

	return toDomainTeam(model), nil
}

func (t *TeamRepository) AddUsersToTeam(ctx context.Context, teamId int, users []domain.User) ([]domain.User, error) {
	const op = "UserRepository.AddUsersToTeam"

	var userIDs []string
	var userNames []string
	var userActives []bool

	for _, user := range users {
		userIDs = append(userIDs, user.Id)
		userNames = append(userNames, user.Name)
		userActives = append(userActives, true)
	}

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

	rows, err := t.Pool.Query(ctx, query, userIDs, userNames, userActives, teamId)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer rows.Close()

	updUsers := make([]domain.User, 0, len(users))
	for rows.Next() {
		var userId string
		var userName string
		var userIsActive bool
		var userTeamId *int
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

func toDomainTeam(model TeamModel) domain.Team {
	return domain.Team{
		Id:   model.Id,
		Name: model.Name,
	}
}

func toTeamModel(team domain.Team) TeamModel {
	return TeamModel{
		Id:   team.Id,
		Name: team.Name,
	}
}
