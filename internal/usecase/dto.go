package usecase

import (
	"avito-internship/internal/domain"
	r "avito-internship/internal/repository"
	"time"
)

type TeamMemberDTO struct {
	Id       string
	Username string
	IsActive bool
}

type TeamDTO struct {
	TeamName string
	Members  []TeamMemberDTO
}

type UserDTO struct {
	Id       string
	Username string
	TeamName string
	IsActive bool
}

type PullRequestDTO struct {
	Id                string
	Name              string
	AuthorId          string
	Status            domain.PRStatus
	AssignedReviewers []string
	CreatedAt         *string
	MergedAt          *string
}

type PullRequestShort struct {
	Id       string
	Name     string
	AuthorId string
	Status   domain.PRStatus
}

type TeamAddReq struct {
	TeamName string
	Members  []TeamMemberDTO
}

type TeamAddRes struct {
	Team TeamDTO
}

type GetTeamQueryReq struct {
	TeamName string
}

type GetTeamRes struct {
	TeamName string
	Members  []TeamMemberDTO
}

type SetIsActiveReq struct {
	UserId   string
	IsActive bool
}

type SetIsActiveRes struct {
	User UserDTO
}

type CreatePullRequestReq struct {
	Id       string
	Name     string
	AuthorId string
}

type CreatePullRequestRes struct {
	PullRequest PullRequestDTO
}

type PullRequestMergeReq struct {
	Id string
}

type PullRequestMergeRes struct {
	PullRequest PullRequestDTO
}

type PullRequestReassignReq struct {
	PullRequestId string
	OldReviewerId string
}

type PullRequestReassignRes struct {
	Pr         PullRequestDTO
	ReplacedBy string
}

type GetReviewQueryReq struct {
	UserID string
}

type GetReviewRes struct {
	UserId       string
	PullRequests []PullRequestShort
}

func NewSetIsActiveRes(id, username, teamName string, isActive bool) SetIsActiveRes {
	return SetIsActiveRes{
		User: UserDTO{
			Id:       id,
			Username: username,
			TeamName: teamName,
			IsActive: isActive,
		},
	}
}

func NewGetReviewRes(userId string, prs r.GetPRByReviewerDTO) GetReviewRes {
	return GetReviewRes{
		UserId:       userId,
		PullRequests: toArrPullRequestShort(prs),
	}
}

func toArrPullRequestShort(prs r.GetPRByReviewerDTO) []PullRequestShort {
	result := make([]PullRequestShort, 0, len(prs.Prs))

	for _, dto := range prs.Prs {
		result = append(result, toPullRequestShort(dto.Pr, dto.StatusName))

	}

	return result
}

// TODO: переделать
func toPullRequestShort(pr domain.PullRequest, statusName domain.PRStatus) PullRequestShort {
	return PullRequestShort{
		Id:       pr.Id,
		Name:     pr.Name,
		AuthorId: pr.AuthorId,
		Status:   statusName,
	}
}

func TeamMemberDTOtoDomainUser(dto TeamMemberDTO) domain.User {
	return domain.User{
		Id:       dto.Id,
		Name:     dto.Username,
		IsActive: dto.IsActive,
	}
}

func NewTeamAddRes(teamDTO TeamDTO) TeamAddRes {
	return TeamAddRes{
		Team: teamDTO,
	}
}

func toArrTeamMemberDTO(u []domain.User) []TeamMemberDTO {
	result := make([]TeamMemberDTO, 0)
	for _, user := range u {
		result = append(result, toTeamMemberDTO(user))
	}

	return result
}

func toTeamMemberDTO(user domain.User) TeamMemberDTO {
	return TeamMemberDTO{
		Id:       user.Id,
		Username: user.Name,
		IsActive: user.IsActive,
	}
}

func NewTeamDTO(teamName string, u []domain.User) TeamDTO {
	return TeamDTO{
		TeamName: teamName,
		Members:  toArrTeamMemberDTO(u),
	}
}

func NewPullRequestDTO(pr domain.PullRequest, reviewers []string, statusName domain.PRStatus) PullRequestDTO {
	createdAt := pr.CreatedAt.Format(time.RFC3339)
	var mergedAt *string
	if pr.MergedAt != nil {
		s := pr.MergedAt.Format(time.RFC3339)
		mergedAt = &s
	}

	result := PullRequestDTO{
		Id:                pr.Id,
		Name:              pr.Name,
		AuthorId:          pr.AuthorId,
		Status:            statusName,
		AssignedReviewers: reviewers,
		CreatedAt:         &createdAt,
		MergedAt:          mergedAt,
	}

	return result
}

func NewCreatePullRequestRes(prDTO PullRequestDTO) CreatePullRequestRes {
	return CreatePullRequestRes{
		PullRequest: prDTO,
	}
}

func NewPullRequestReassignRes(pr PullRequestDTO, replacedBy string) PullRequestReassignRes {
	return PullRequestReassignRes{
		Pr:         pr,
		ReplacedBy: replacedBy,
	}
}

func NewPullRequestMergeRes(pr PullRequestDTO) PullRequestMergeRes {
	return PullRequestMergeRes{
		PullRequest: pr,
	}
}

func NewGetTeamRes(teamDTO TeamDTO) GetTeamRes {
	return GetTeamRes{
		TeamName: teamDTO.TeamName,
		Members:  teamDTO.Members,
	}
}
