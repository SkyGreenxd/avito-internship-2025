package usecase

import (
	"avito-internship/internal/domain"
	r "avito-internship/internal/repository"
	"avito-internship/pkg/e"
	"avito-internship/pkg/transaction"
	"context"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
)

type TeamUseCase struct {
	teamRepo     r.TeamRepository
	userRepo     r.UserRepository
	prRepo       r.PullRequestRepository
	statusRepo   r.StatusRepository
	reviewerRepo r.PrReviewerRepository
	dbPool       transaction.Transactional
}

func NewTeamUseCase(teamRepo r.TeamRepository, userRepo r.UserRepository,
	prRepo r.PullRequestRepository, statusRepo r.StatusRepository,
	dbPool transaction.Transactional, reviewerRepo r.PrReviewerRepository) *TeamUseCase {
	return &TeamUseCase{
		teamRepo:     teamRepo,
		userRepo:     userRepo,
		prRepo:       prRepo,
		statusRepo:   statusRepo,
		reviewerRepo: reviewerRepo,
		dbPool:       dbPool,
	}
}

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

func (t *TeamUseCase) GetTeam(ctx context.Context, teamName string) (GetTeamRes, error) {
	const op = "TeamUseCase.GetTeam"

	members, err := t.teamRepo.GetMembersByTeamNameWithUsers(ctx, teamName)
	if err != nil {
		return GetTeamRes{}, e.Wrap(op, err)
	}

	teamDTO := NewTeamDTO(teamName, members)
	return NewGetTeamRes(teamDTO), nil
}

func (t *TeamUseCase) DeactivateMembers(ctx context.Context, req DeactivateMembersReq) (DeactivateMembersRes, error) {
	const op = "TeamUseCase.DeactivateMembers"

	if len(req.Members) == 0 {
		return DeactivateMembersRes{}, e.Wrap(op, e.ErrEmptyMembers)
	}

	allMembers, err := t.teamRepo.GetMembersByTeamNameWithUsers(ctx, req.TeamName)
	if err != nil {
		return DeactivateMembersRes{}, e.Wrap(op, err)
	}

	teamMembersSet := make(map[string]struct{}, len(allMembers))
	for _, m := range allMembers {
		teamMembersSet[m.Id] = struct{}{}
	}

	for _, id := range req.Members {
		if _, ok := teamMembersSet[id]; !ok {
			return DeactivateMembersRes{}, e.Wrap(op, e.ErrInvalidMember)
		}
	}

	status, err := t.statusRepo.GetByName(ctx, string(domain.OPEN))
	if err != nil {
		return DeactivateMembersRes{}, e.Wrap(op, err)
	}

	idSet := make(map[string]struct{}, len(req.Members))
	deactivateMembersIds := make([]string, 0, len(req.Members))
	for _, id := range req.Members {
		idSet[id] = struct{}{}
		deactivateMembersIds = append(deactivateMembersIds, id)
	}

	globalCandidatePool := make(map[string]struct{})
	for _, member := range allMembers {
		if member.IsActive {
			if _, deactivating := idSet[member.Id]; !deactivating {
				globalCandidatePool[member.Id] = struct{}{}
			}
		}
	}

	if len(globalCandidatePool) == 0 {
		return DeactivateMembersRes{}, e.Wrap(op, e.ErrPrNoCandidate)
	}

	ctx, tx, err := transaction.NewTransaction(ctx, pgx.TxOptions{}, t.dbPool)
	if err != nil {
		return DeactivateMembersRes{}, e.Wrap(op, err)
	}
	defer tx.Rollback(ctx)
	ctx = context.WithValue(ctx, "tx", tx.Transaction())

	prMap, err := t.prRepo.GetOpenPRsByReviewerIDs(ctx, deactivateMembersIds, status.Id)
	if err != nil {
		return DeactivateMembersRes{}, e.Wrap(op, err)
	}

	rd := rand.New(rand.NewSource(time.Now().UnixNano()))

	prUpdates := make(map[string][]string)
	prChanges := make(map[string]r.PrReviewerChange)

	for prId, data := range prMap {
		pr := data.Pr
		allReviewersForThisPR := data.ReviewersIds

		reviewersToReplace := make([]string, 0)
		activeReviewersOnPR := make([]string, 0)

		for _, reviewerId := range allReviewersForThisPR {
			if _, isDeactivating := idSet[reviewerId]; isDeactivating {
				reviewersToReplace = append(reviewersToReplace, reviewerId)
			} else {
				activeReviewersOnPR = append(activeReviewersOnPR, reviewerId)
			}
		}

		if len(reviewersToReplace) == 0 {
			continue
		}

		prCandidatePool := make([]string, 0)
		for id := range globalCandidatePool {
			if id != pr.AuthorId {
				prCandidatePool = append(prCandidatePool, id)
			}
		}

		existingSet := make(map[string]struct{})
		for _, r := range allReviewersForThisPR {
			existingSet[r] = struct{}{}
		}

		cleanCandidates := make([]string, 0)
		for _, id := range prCandidatePool {
			if _, exists := existingSet[id]; !exists {
				cleanCandidates = append(cleanCandidates, id)
			}
		}

		if len(cleanCandidates) < len(reviewersToReplace) {
			return DeactivateMembersRes{}, e.Wrap(op, e.ErrPrNoCandidate)
		}

		rd.Shuffle(len(cleanCandidates), func(i, j int) {
			cleanCandidates[i], cleanCandidates[j] = cleanCandidates[j], cleanCandidates[i]
		})
		newReviewers := cleanCandidates[:len(reviewersToReplace)]

		prUpdates[prId] = append(activeReviewersOnPR, newReviewers...)
		prChanges[prId] = r.PrReviewerChange{
			ToAdd:    newReviewers,
			ToRemove: reviewersToReplace,
		}
	}

	if len(prChanges) > 0 {
		err = t.reviewerRepo.UpdateReviewers(ctx, prChanges)
		if err != nil {
			return DeactivateMembersRes{}, e.Wrap(op, err)
		}
	}

	updUsers, err := t.userRepo.DeactivateUsers(ctx, deactivateMembersIds)
	if err != nil {
		return DeactivateMembersRes{}, e.Wrap(op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return DeactivateMembersRes{}, e.Wrap(op, err)
	}

	updatedPRs := make([]PullRequestDTO, 0, len(prUpdates))
	for prId, reviewers := range prUpdates {
		data, exists := prMap[prId]
		if !exists {
			continue
		}

		var createdAt *string
		if !data.Pr.CreatedAt.IsZero() {
			t := data.Pr.CreatedAt.Format(time.RFC3339)
			createdAt = &t
		}

		var mergedAt *string
		if data.Pr.MergedAt != nil && !data.Pr.MergedAt.IsZero() {
			t := data.Pr.MergedAt.Format(time.RFC3339)
			mergedAt = &t
		}

		dto := PullRequestDTO{
			Id:                data.Pr.Id,
			Name:              data.Pr.Name,
			AuthorId:          data.Pr.AuthorId,
			Status:            domain.PRStatus(data.StatusName),
			AssignedReviewers: reviewers,
			CreatedAt:         createdAt,
			MergedAt:          mergedAt,
		}

		updatedPRs = append(updatedPRs, dto)
	}

	return NewDeactivateMembersRes(req.TeamName, toArrTeamMemberDTO(updUsers), updatedPRs), nil
}
