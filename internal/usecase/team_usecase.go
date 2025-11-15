package usecase

import (
	"avito-internship/internal/domain"
	r "avito-internship/internal/repository"
	"avito-internship/pkg/e"
	"avito-internship/pkg/transaction"
	"context"

	"github.com/jackc/pgx/v5"
)

type TeamUseCase struct {
	teamRepo r.TeamRepository
	userRepo r.UserRepository
	dbPool   transaction.Transactional
}

func NewTeamUseCase(teamRepo r.TeamRepository, userRepo r.UserRepository, dbPool transaction.Transactional) *TeamUseCase {
	return &TeamUseCase{
		teamRepo: teamRepo,
		userRepo: userRepo,
		dbPool:   dbPool,
	}
}

// AddTeam создает новую команду с участниками
func (t *TeamUseCase) AddTeam(ctx context.Context, req TeamAddReq) (TeamAddRes, error) {
	const op = "TeamUseCase.AddTeam"

	if len(req.Members) == 0 {
		return TeamAddRes{}, e.Wrap(op, e.ErrEmptyMembers)
	}

	ctx, tx, err := transaction.NewTransaction(ctx, pgx.TxOptions{}, t.dbPool)
	if err != nil {
		return TeamAddRes{}, e.Wrap(op, err)
	}
	defer tx.Rollback(ctx)
	ctx = context.WithValue(ctx, "tx", tx.Transaction())

	team := domain.NewTeam(req.TeamName)
	newTeam, err := t.teamRepo.Create(ctx, team)
	if err != nil {
		return TeamAddRes{}, e.Wrap(op, err)
	}

	users := make([]domain.User, 0, len(req.Members))
	for _, member := range req.Members {
		users = append(users, TeamMemberDTOtoDomainUser(member))
	}

	members, err := t.userRepo.AddUsersToTeam(ctx, newTeam.Id, users)
	if err != nil {
		return TeamAddRes{}, e.Wrap(op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return TeamAddRes{}, e.Wrap(op, err)
	}

	teamDTO := NewTeamDTO(newTeam.Name, members)
	return NewTeamAddRes(teamDTO), nil
}

// GetTeam возвращает команду с ее участниками
func (t *TeamUseCase) GetTeam(ctx context.Context, teamName string) (GetTeamRes, error) {
	const op = "TeamUseCase.GetTeam"

	members, err := t.teamRepo.GetMembersByTeamNameWithUsers(ctx, teamName)
	if err != nil {
		return GetTeamRes{}, e.Wrap(op, err)
	}

	teamDTO := NewTeamDTO(teamName, members)
	return NewGetTeamRes(teamDTO), nil
}
