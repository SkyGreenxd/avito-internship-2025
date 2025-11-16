package v1

import "avito-internship/internal/usecase"

func toUseCaseSetIsActiveReq(req SetIsActiveReq) usecase.SetIsActiveReq {
	return usecase.SetIsActiveReq{
		UserId:   req.UserId,
		IsActive: *req.IsActive,
	}
}

func toDeliverySetIsActiveRes(req usecase.SetIsActiveRes) SetIsActiveRes {
	return SetIsActiveRes{
		User: toDeliveryUserDTO(req.User),
	}
}

func toDeliveryUserDTO(u usecase.UserDTO) UserDTO {
	return UserDTO{
		Id:       u.Id,
		Username: u.Username,
		TeamName: u.TeamName,
		IsActive: u.IsActive,
	}
}

func toDeliveryGetReviewRes(req usecase.GetReviewRes) GetReviewRes {
	return GetReviewRes{
		UserId:       req.UserId,
		PullRequests: toArrDeliveryPullRequestShort(req.PullRequests),
	}
}

func toArrDeliveryPullRequestShort(prs []usecase.PullRequestShort) []PullRequestShort {
	result := make([]PullRequestShort, 0, len(prs))
	for _, pr := range prs {
		result = append(result, toDeliveryPullRequestShort(pr))
	}

	return result
}

func toDeliveryPullRequestShort(pr usecase.PullRequestShort) PullRequestShort {
	return PullRequestShort{
		Id:       pr.Id,
		Name:     pr.Name,
		AuthorId: pr.AuthorId,
		Status:   pr.Status,
	}
}

func toDeliveryGetTeamRes(res usecase.GetTeamRes) GetTeamRes {
	return GetTeamRes{
		TeamName: res.TeamName,
		Members:  toArrTeamMemberDTO(res.Members),
	}
}

func toDeliveryTeamAddRes(res usecase.TeamAddRes) TeamAddRes {
	return TeamAddRes{
		Team: toDeliveryTeamDTO(res.Team),
	}
}

func toUseCaseTeamAddReq(req TeamAddReq) usecase.TeamAddReq {
	return usecase.TeamAddReq{
		TeamName: req.TeamName,
		Members:  toArrUseCaseTeamMemberDTO(req.Members),
	}
}

func toArrUseCaseTeamMemberDTO(ms []TeamMemberDTO) []usecase.TeamMemberDTO {
	result := make([]usecase.TeamMemberDTO, 0, len(ms))
	for _, m := range ms {
		result = append(result, toUseCaseTeamMemberDTO(m))
	}

	return result
}

func toUseCaseTeamMemberDTO(m TeamMemberDTO) usecase.TeamMemberDTO {
	return usecase.TeamMemberDTO{
		Id:       m.Id,
		Username: m.Username,
		IsActive: *m.IsActive,
	}
}

func toDeliveryTeamDTO(res usecase.TeamDTO) TeamDTO {
	return TeamDTO{
		TeamName: res.TeamName,
		Members:  toArrTeamMemberDTO(res.Members),
	}
}

func toArrTeamMemberDTO(ms []usecase.TeamMemberDTO) []TeamMemberDTO {
	result := make([]TeamMemberDTO, 0, len(ms))
	for _, m := range ms {
		result = append(result, toDeliveryTeamMemberDTO(m))
	}
	return result
}

func toDeliveryTeamMemberDTO(m usecase.TeamMemberDTO) TeamMemberDTO {
	return TeamMemberDTO{
		Id:       m.Id,
		Username: m.Username,
		IsActive: &m.IsActive,
	}
}

func toUseCaseCreatePullRequestReq(req CreatePullRequestReq) usecase.CreatePullRequestReq {
	return usecase.CreatePullRequestReq{
		Id:       req.Id,
		Name:     req.Name,
		AuthorId: req.AuthorId,
	}
}

func toDeliveryCreatePullRequestRes(res usecase.CreatePullRequestRes) CreatePullRequestRes {
	return CreatePullRequestRes{
		PullRequest: toCreatePullRequestDTO(res.PullRequest),
	}
}

func toCreatePullRequestDTO(pr usecase.PullRequestDTO) CreatePullRequestDTO {
	return CreatePullRequestDTO{
		Id:                pr.Id,
		Name:              pr.Name,
		AuthorId:          pr.AuthorId,
		Status:            pr.Status,
		AssignedReviewers: pr.AssignedReviewers,
	}
}

func toDeliveryPullRequestDTO(pr usecase.PullRequestDTO) PullRequestDTO {
	return PullRequestDTO{
		Id:                pr.Id,
		Name:              pr.Name,
		AuthorId:          pr.AuthorId,
		Status:            pr.Status,
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func toUseCasePullRequestMergeReq(req PullRequestMergeReq) usecase.PullRequestMergeReq {
	return usecase.PullRequestMergeReq{
		Id: req.Id,
	}
}

func toDeliveryPullRequestMergeRes(res usecase.PullRequestMergeRes) PullRequestMergeRes {
	return PullRequestMergeRes{
		PullRequest: NewMergePullRequestDTO(res.PullRequest),
	}
}

func NewMergePullRequestDTO(pr usecase.PullRequestDTO) MergePullRequestDTO {
	return MergePullRequestDTO{
		Id:                pr.Id,
		Name:              pr.Name,
		AuthorId:          pr.AuthorId,
		Status:            pr.Status,
		AssignedReviewers: pr.AssignedReviewers,
		MergedAt:          pr.MergedAt,
	}
}

func toUseCasePullRequestReassignReq(req PullRequestReassignReq) usecase.PullRequestReassignReq {
	return usecase.PullRequestReassignReq{
		PullRequestId: req.PullRequestId,
		OldReviewerId: req.OldReviewerId,
	}
}

func toDeliveryPullRequestReassignRes(res usecase.PullRequestReassignRes) PullRequestReassignRes {
	return PullRequestReassignRes{
		Pr:         toDeliveryPullRequestDTO(res.Pr),
		ReplacedBy: res.ReplacedBy,
	}
}

func toUseCaseDeactivateMembers(req DeactivateMembersReq) usecase.DeactivateMembersReq {
	return usecase.DeactivateMembersReq{
		TeamName: req.TeamName,
		Members:  req.Members,
	}
}

func toDeliveryDeactivateMembers(res usecase.DeactivateMembersRes) DeactivateMembersRes {
	return DeactivateMembersRes{
		TeamName:           res.TeamName,
		DeactivatedMembers: toArrTeamMemberDTO(res.DeactivatedMembers),
		UpdPrs:             toArrDeliveryPullRequestDTO(res.UpdPrs),
	}
}

func toArrDeliveryPullRequestDTO(pr []usecase.PullRequestDTO) []PullRequestDTO {
	result := make([]PullRequestDTO, 0, len(pr))
	for _, pr := range pr {
		result = append(result, toDeliveryPullRequestDTO(pr))
	}
	return result
}
