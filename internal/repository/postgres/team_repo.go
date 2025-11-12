package postgres

import (
	"avito-internship/internal/domain"
	"avito-internship/pkg/e"
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	Pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{Pool: pool}
}

func (t *TeamRepository) Create(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	const op = "TeamRepository.Create"

	model := toTeamModel(team)
	// TODO: перенести в UseCase
	model.Id = uuid.NewString()

	queryBuilder := sq.Insert("teams").
		Columns("id", "name").
		Values(model.Id, model.Name).
		Suffix("RETURNING id, name")

	query, args, err := queryBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	err = t.Pool.QueryRow(ctx, query, args...).Scan(&model.Id, &model.Name)
	if err = postgresDuplicate(err, e.ErrTeamIsExists); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTeam(model), nil
}

func (t *TeamRepository) GetByTeamNameWithUsers(ctx context.Context, teamName string) (*domain.TeamWithUsers, error) {
	return t.getByTeamNameWithUsers(ctx, teamName, false)
}

func (t *TeamRepository) GetByTeamNameWithActiveUsers(ctx context.Context, teamName string) (*domain.TeamWithUsers, error) {
	return t.getByTeamNameWithUsers(ctx, teamName, true)
}

func (t *TeamRepository) getByTeamNameWithUsers(ctx context.Context, teamName string, onlyActive bool) (*domain.TeamWithUsers, error) {
	const op = "TeamRepository.getByTeamNameWithUsers"

	// Получаем команду
	tBuilder := sq.Select("id", "name").
		From("teams").
		Where(sq.Eq{"name": teamName})
	tQuery, tArgs, err := tBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	var teamModel TeamModel
	err = t.Pool.QueryRow(ctx, tQuery, tArgs...).Scan(&teamModel.Id, &teamModel.Name)
	if err := checkGetQueryResult(err, e.ErrTeamNotFound); err != nil {
		return nil, e.Wrap(op, err)
	}

	// Получаем пользователей
	uBuilder := sq.Select("id", "name", "team_id").From("users").Where(sq.Eq{"team_id": teamModel.Id})
	if onlyActive {
		uBuilder = uBuilder.Where(sq.Eq{"is_active": true})
	}

	uQuery, uArgs, err := uBuilder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, e.Wrap(op, err)
	}

	rows, err := t.Pool.Query(ctx, uQuery, uArgs...)
	if err != nil {
		return nil, e.Wrap(op, err)
	}
	defer rows.Close()

	usersModel := make([]UserModel, 0)
	for rows.Next() {
		var userModel UserModel
		if err := rows.Scan(&userModel.Id, &userModel.Name, &userModel.TeamId); err != nil {
			return nil, e.Wrap(op, err)
		}
		usersModel = append(usersModel, userModel)
	}
	if err := rows.Err(); err != nil {
		return nil, e.Wrap(op, err)
	}

	return toDomainTeamWithUsers(&teamModel, usersModel), nil
}

func toTeamModel(t *domain.Team) *TeamModel {
	return &TeamModel{
		Id:   t.Id,
		Name: t.Name,
	}
}

func toDomainTeam(t *TeamModel) *domain.Team {
	return &domain.Team{
		Id:   t.Id,
		Name: t.Name,
	}
}

func toDomainTeamWithUsers(t *TeamModel, u []UserModel) *domain.TeamWithUsers {
	team := toDomainTeam(t)

	users := make([]domain.User, 0, len(u))
	for _, user := range u {
		users = append(users, *toDomainUser(&user))
	}

	return &domain.TeamWithUsers{
		Team:  team,
		Users: users,
	}
}
