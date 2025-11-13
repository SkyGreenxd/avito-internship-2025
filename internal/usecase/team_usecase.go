package usecase

import (
	"avito-internship/internal/domain"
	r "avito-internship/internal/repository"
	"avito-internship/pkg/e"
	"context"
)

type TeamUseCase struct {
	teamRepo r.TeamRepository
}

func NewTeamUseCase(teamRepo r.TeamRepository) *TeamUseCase {
	return &TeamUseCase{
		teamRepo: teamRepo,
	}
}

// AddTeam создает новую команду с участниками
func (t *TeamUseCase) AddTeam(ctx context.Context, req TeamDTO) (TeamAddRes, error) {
	const op = "TeamUseCase.AddTeam"

	team := domain.NewTeam(req.TeamName)
	newTeam, err := t.teamRepo.Create(ctx, *team)
	if err != nil {
		return TeamAddRes{}, e.Wrap(op, err)
	}

	users := make([]domain.User, 0, len(req.Members))
	for _, member := range req.Members {
		users = append(users, TeamMemberDTOtoDomainUser(member))
	}

	members, err := t.teamRepo.AddUsersToTeam(ctx, newTeam.Id, users)
	if err != nil {
		return TeamAddRes{}, e.Wrap(op, err)
	}

	teamDTO := NewTeamDTO(newTeam.Name, members)
	return NewTeamAddRes(teamDTO), nil
}

// GetTeam возвращает команду с ее участниками
func (t *TeamUseCase) GetTeam(ctx context.Context, teamName string) (TeamDTO, error) {
	const op = "TeamUseCase.GetTeam"

	members, err := t.teamRepo.GetMembersByTeamNameWithUsers(ctx, teamName)
	if err != nil {
		return TeamDTO{}, e.Wrap(op, err)
	}

	teamDTO := NewTeamDTO(teamName, members)
	return teamDTO, nil
}
