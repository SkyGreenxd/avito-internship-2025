package usecase

import (
	"avito-internship/internal/domain"
	r "avito-internship/internal/repository"
	repoMocks "avito-internship/internal/repository/mocks"
	"avito-internship/pkg/e"
	trMock "avito-internship/pkg/transaction/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPullRequestUseCase_PullRequestCreate(t *testing.T) {
	fixedTime := time.Date(2025, 11, 15, 23, 54, 48, 0, time.UTC)
	createdAtStr := fixedTime.Format(time.RFC3339)

	tests := []struct {
		name              string
		req               CreatePullRequestReq
		statusRepoSetup   func(*repoMocks.MockStatusRepository)
		prRepoSetup       func(*repoMocks.MockPullRequestRepository)
		userRepoSetup     func(*repoMocks.MockUserRepository)
		reviewerRepoSetup func(repository *repoMocks.MockPrReviewerRepository)
		expectedRes       CreatePullRequestRes
		expectedErr       error
	}{
		{
			name: "success",
			req: CreatePullRequestReq{
				Id:       "pr-1001",
				Name:     "Test PR",
				AuthorId: "u1",
			},
			statusRepoSetup: func(repo *repoMocks.MockStatusRepository) {
				repo.EXPECT().GetByName(gomock.Any(), "OPEN").
					Return(domain.Status{1, "OPEN"}, nil)
			},
			prRepoSetup: func(repo *repoMocks.MockPullRequestRepository) {
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error) {
						pr.CreatedAt = fixedTime
						return pr, nil
					},
				)
			},
			userRepoSetup: func(repo *repoMocks.MockUserRepository) {
				repo.EXPECT().
					GetReviewCandidates(gomock.Any(), "u1", 2).
					Return([]domain.User{
						{
							Id:       "u2",
							Name:     "Test PR",
							IsActive: true,
							TeamId:   1,
						},
						{
							Id:       "u3",
							Name:     "Test PR_2",
							IsActive: true,
							TeamId:   1,
						},
					}, nil)
			},
			reviewerRepoSetup: func(repository *repoMocks.MockPrReviewerRepository) {
				repository.EXPECT().AddReviewers(gomock.Any(), "pr-1001", []string{"u2", "u3"}).
					Return(nil)
			},
			expectedRes: CreatePullRequestRes{
				PullRequest: PullRequestDTO{
					Id:                "pr-1001",
					Name:              "Test PR",
					AuthorId:          "u1",
					Status:            domain.OPEN,
					AssignedReviewers: []string{"u2", "u3"},
					CreatedAt:         &createdAtStr,
					MergedAt:          nil,
				},
			},
			expectedErr: nil,
		},
		{
			name: "pr already exists",
			req: CreatePullRequestReq{
				Id:       "pr-1001",
				Name:     "Existing PR",
				AuthorId: "u1",
			},
			statusRepoSetup: func(repo *repoMocks.MockStatusRepository) {
				repo.EXPECT().GetByName(gomock.Any(), "OPEN").
					Return(domain.Status{1, "OPEN"}, nil)
			},
			prRepoSetup: func(repo *repoMocks.MockPullRequestRepository) {
				repo.EXPECT().Create(gomock.Any(), gomock.Any()).
					Return(domain.PullRequest{}, e.ErrPRIsExists)
			},
			userRepoSetup: func(repo *repoMocks.MockUserRepository) {
				repo.EXPECT().GetReviewCandidates(gomock.Any(), "u1", 2).Return([]domain.User{}, nil)
			},
			reviewerRepoSetup: func(repo *repoMocks.MockPrReviewerRepository) {},
			expectedRes:       CreatePullRequestRes{},
			expectedErr:       e.ErrPRIsExists,
		},
		{
			name: "user not found",
			req: CreatePullRequestReq{
				Id:       "pr-1001",
				Name:     "Existing PR",
				AuthorId: "u888",
			},
			statusRepoSetup: func(repo *repoMocks.MockStatusRepository) {},
			prRepoSetup:     func(repo *repoMocks.MockPullRequestRepository) {},
			userRepoSetup: func(repo *repoMocks.MockUserRepository) {
				repo.EXPECT().GetReviewCandidates(gomock.Any(), "u888", 2).Return([]domain.User{}, e.ErrUserNotFound)
			},
			reviewerRepoSetup: func(repo *repoMocks.MockPrReviewerRepository) {},
			expectedRes:       CreatePullRequestRes{},
			expectedErr:       e.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			statusRepo := repoMocks.NewMockStatusRepository(ctrl)
			userRepo := repoMocks.NewMockUserRepository(ctrl)
			prRepo := repoMocks.NewMockPullRequestRepository(ctrl)
			reviewerRepo := repoMocks.NewMockPrReviewerRepository(ctrl)

			mockTx := trMock.NewMockTx(ctrl)
			mockTx.EXPECT().Commit(gomock.Any()).Return(nil).AnyTimes()
			mockTx.EXPECT().Rollback(gomock.Any()).Return(nil).AnyTimes()

			mockTxPool := trMock.NewMockTransactional(ctrl)
			mockTxPool.EXPECT().
				BeginTx(gomock.Any(), gomock.Any()).
				Return(mockTx, nil).
				AnyTimes()

			prUC := NewPullRequestUseCase(prRepo, reviewerRepo, userRepo, statusRepo, mockTxPool)

			tt.statusRepoSetup(statusRepo)
			tt.prRepoSetup(prRepo)
			tt.userRepoSetup(userRepo)
			tt.reviewerRepoSetup(reviewerRepo)

			res, err := prUC.PullRequestCreate(context.Background(), tt.req)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error: got %v, want %v", err, tt.expectedErr)
			}

			require.Equal(t, tt.expectedRes.PullRequest.Id, res.PullRequest.Id)
			require.Equal(t, tt.expectedRes.PullRequest.Name, res.PullRequest.Name)
			require.Equal(t, tt.expectedRes.PullRequest.AuthorId, res.PullRequest.AuthorId)
			require.Equal(t, tt.expectedRes.PullRequest.Status, res.PullRequest.Status)
			require.ElementsMatch(t, tt.expectedRes.PullRequest.AssignedReviewers, res.PullRequest.AssignedReviewers)
			require.Equal(t, tt.expectedRes.PullRequest.MergedAt, res.PullRequest.MergedAt)
		})
	}
}

func TestPullRequestUseCase_PullRequestMerge(t *testing.T) {
	fixedTime := time.Date(2025, 11, 15, 23, 54, 48, 0, time.UTC)
	mergedAtStr := fixedTime.Format(time.RFC3339)

	tests := []struct {
		name            string
		req             PullRequestMergeReq
		statusRepoSetup func(*repoMocks.MockStatusRepository)
		prRepoSetup     func(*repoMocks.MockPullRequestRepository)
		expectedRes     PullRequestMergeRes
		expectedErr     error
	}{
		{
			name: "success",
			req: PullRequestMergeReq{
				Id: "pr-1001",
			},
			statusRepoSetup: func(repo *repoMocks.MockStatusRepository) {
				repo.EXPECT().GetByName(gomock.Any(), "MERGED").
					Return(domain.Status{Id: 2, Name: "MERGED"}, nil)
			},
			prRepoSetup: func(repo *repoMocks.MockPullRequestRepository) {
				repo.EXPECT().SetMergedStatus(gomock.Any(), 2, "pr-1001").
					DoAndReturn(func(ctx context.Context, statusId int, prId string) (r.SetMergedStatusDTO, error) {
						pr := domain.PullRequest{
							Id:                prId,
							Name:              "Test PR",
							AuthorId:          "u1",
							StatusId:          statusId,
							NeedMoreReviewers: false,
							CreatedAt:         fixedTime,
							MergedAt:          &fixedTime,
						}
						return r.SetMergedStatusDTO{
							Pr:           pr,
							ReviewersIds: []string{"u2", "u3"},
						}, nil
					})
			},
			expectedRes: PullRequestMergeRes{
				PullRequest: PullRequestDTO{
					Id:                "pr-1001",
					Name:              "Test PR",
					AuthorId:          "u1",
					Status:            domain.MERGED,
					AssignedReviewers: []string{"u2", "u3"},
					CreatedAt:         &mergedAtStr,
					MergedAt:          &mergedAtStr,
				},
			},
			expectedErr: nil,
		},
		{
			name: "pr not found",
			req: PullRequestMergeReq{
				Id: "pr-9999",
			},
			statusRepoSetup: func(repo *repoMocks.MockStatusRepository) {
				repo.EXPECT().GetByName(gomock.Any(), "MERGED").
					Return(domain.Status{Id: 2, Name: "MERGED"}, nil)
			},
			prRepoSetup: func(repo *repoMocks.MockPullRequestRepository) {
				repo.EXPECT().SetMergedStatus(gomock.Any(), 2, "pr-9999").
					Return(r.SetMergedStatusDTO{}, e.ErrPRNotFound)
			},
			expectedRes: PullRequestMergeRes{},
			expectedErr: e.ErrPRNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			statusRepo := repoMocks.NewMockStatusRepository(ctrl)
			prRepo := repoMocks.NewMockPullRequestRepository(ctrl)

			tt.statusRepoSetup(statusRepo)
			tt.prRepoSetup(prRepo)

			prUC := NewPullRequestUseCase(prRepo, nil, nil, statusRepo, nil)

			res, err := prUC.PullRequestMerge(context.Background(), tt.req)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error: got %v, want %v", err, tt.expectedErr)
			}
			require.Equal(t, tt.expectedRes.PullRequest.Id, res.PullRequest.Id)
			require.Equal(t, tt.expectedRes.PullRequest.Name, res.PullRequest.Name)
			require.Equal(t, tt.expectedRes.PullRequest.AuthorId, res.PullRequest.AuthorId)
			require.Equal(t, tt.expectedRes.PullRequest.Status, res.PullRequest.Status)
			require.ElementsMatch(t, tt.expectedRes.PullRequest.AssignedReviewers, res.PullRequest.AssignedReviewers)
			require.Equal(t, tt.expectedRes.PullRequest.CreatedAt, res.PullRequest.CreatedAt)
			require.Equal(t, tt.expectedRes.PullRequest.MergedAt, res.PullRequest.MergedAt)
		})
	}
}

func TestPullRequestUseCase_ReviewerReassign(t *testing.T) {
	tests := []struct {
		name              string
		input             PullRequestReassignReq
		userRepoSetup     func(*repoMocks.MockUserRepository)
		prRepoSetup       func(*repoMocks.MockPullRequestRepository)
		reviewerRepoSetup func(repository *repoMocks.MockPrReviewerRepository)
		expectedRes       PullRequestReassignRes
		expectedErr       error
	}{
		{
			name: "success",
			input: PullRequestReassignReq{
				PullRequestId: "pr-1001",
				OldReviewerId: "u3",
			},
			userRepoSetup: func(repo *repoMocks.MockUserRepository) {
				repo.EXPECT().GetById(gomock.Any(), "u3").Return(domain.User{
					Id:       "u3",
					Name:     "test",
					IsActive: true,
					TeamId:   1,
				}, nil)

				repo.EXPECT().
					GetReassignCandidates(gomock.Any(), "u1", gomock.Any(), gomock.Any()).
					Return([]domain.User{
						{Id: "u4", Name: "newReviewer", IsActive: true, TeamId: 1},
					}, nil)
			},
			prRepoSetup: func(repo *repoMocks.MockPullRequestRepository) {
				repo.EXPECT().
					GetByPrIdWithReviewersIds(gomock.Any(), "pr-1001").
					Return(r.GetByPrIdWithReviewersIdsDTO{
						Pr: domain.PullRequest{
							Id:                "pr-1001",
							Name:              "My PR",
							AuthorId:          "u1",
							StatusId:          1,
							NeedMoreReviewers: false,
							CreatedAt:         time.Now(),
							MergedAt:          nil,
						},
						ReviewersIds: []string{"u2", "u3"},
						StatusName:   domain.OPEN,
					}, nil)
			},
			reviewerRepoSetup: func(repo *repoMocks.MockPrReviewerRepository) {
				repo.EXPECT().
					UpdateReviewer(gomock.Any(), "u3", "u4", "pr-1001").
					Return("u4", nil)
			},
			expectedRes: PullRequestReassignRes{
				Pr: PullRequestDTO{
					Id:                "pr-1001",
					Name:              "My PR",
					AuthorId:          "u1",
					Status:            domain.OPEN,
					AssignedReviewers: []string{"u2", "u4"},
				},
				ReplacedBy: "u4",
			},
			expectedErr: nil,
		},
		{
			name: "user not found",
			input: PullRequestReassignReq{
				PullRequestId: "pr-1001",
				OldReviewerId: "u999",
			},
			userRepoSetup: func(repo *repoMocks.MockUserRepository) {
				repo.EXPECT().
					GetById(gomock.Any(), "u999").
					Return(domain.User{}, e.ErrUserNotFound)
			},
			prRepoSetup:       func(repo *repoMocks.MockPullRequestRepository) {},
			reviewerRepoSetup: func(repo *repoMocks.MockPrReviewerRepository) {},
			expectedRes:       PullRequestReassignRes{},
			expectedErr:       e.ErrUserNotFound,
		},
		{
			name: "pr not found",
			input: PullRequestReassignReq{
				PullRequestId: "pr-404",
				OldReviewerId: "u3",
			},
			userRepoSetup: func(repo *repoMocks.MockUserRepository) {
				repo.EXPECT().
					GetById(gomock.Any(), "u3").
					Return(domain.User{
						Id:       "u3",
						Name:     "test",
						IsActive: true,
						TeamId:   1,
					}, nil)
			},
			prRepoSetup: func(repo *repoMocks.MockPullRequestRepository) {
				repo.EXPECT().
					GetByPrIdWithReviewersIds(gomock.Any(), "pr-404").
					Return(r.GetByPrIdWithReviewersIdsDTO{}, e.ErrPRNotFound)
			},
			reviewerRepoSetup: func(repo *repoMocks.MockPrReviewerRepository) {},
			expectedRes:       PullRequestReassignRes{},
			expectedErr:       e.ErrPRNotFound,
		},
		{
			name: "error_pr_merged",
			input: PullRequestReassignReq{
				PullRequestId: "pr-1002",
				OldReviewerId: "u3",
			},
			userRepoSetup: func(repo *repoMocks.MockUserRepository) {
				repo.EXPECT().GetById(gomock.Any(), "u3").Return(domain.User{
					Id:       "u3",
					Name:     "test",
					IsActive: true,
					TeamId:   1,
				}, nil)
			},
			prRepoSetup: func(repo *repoMocks.MockPullRequestRepository) {
				repo.EXPECT().
					GetByPrIdWithReviewersIds(gomock.Any(), "pr-1002").
					Return(r.GetByPrIdWithReviewersIdsDTO{
						Pr: domain.PullRequest{
							Id:                "pr-1002",
							Name:              "Merged PR",
							AuthorId:          "u1",
							StatusId:          2,
							NeedMoreReviewers: false,
							CreatedAt:         time.Now(),
							MergedAt:          ptr(time.Now()),
						},
						ReviewersIds: []string{"u2", "u3"},
						StatusName:   domain.MERGED,
					}, nil)
			},
			reviewerRepoSetup: func(repo *repoMocks.MockPrReviewerRepository) {},
			expectedRes:       PullRequestReassignRes{},
			expectedErr:       e.ErrPrMerged,
		},
		{
			name: "error_reviewer_not_assigned",
			input: PullRequestReassignReq{
				PullRequestId: "pr-1003",
				OldReviewerId: "u3",
			},
			userRepoSetup: func(repo *repoMocks.MockUserRepository) {
				repo.EXPECT().GetById(gomock.Any(), "u3").Return(domain.User{
					Id:       "u3",
					Name:     "test",
					IsActive: true,
					TeamId:   1,
				}, nil)
			},
			prRepoSetup: func(repo *repoMocks.MockPullRequestRepository) {
				repo.EXPECT().
					GetByPrIdWithReviewersIds(gomock.Any(), "pr-1003").
					Return(r.GetByPrIdWithReviewersIdsDTO{
						Pr: domain.PullRequest{
							Id:                "pr-1003",
							Name:              "No such reviewer",
							AuthorId:          "u1",
							StatusId:          1,
							NeedMoreReviewers: false,
							CreatedAt:         time.Now(),
							MergedAt:          nil,
						},
						ReviewersIds: []string{"u2", "u4"},
						StatusName:   domain.OPEN,
					}, nil)
			},
			reviewerRepoSetup: func(repo *repoMocks.MockPrReviewerRepository) {},
			expectedRes:       PullRequestReassignRes{},
			expectedErr:       e.ErrPrReviewerNotAssigned,
		},
		{
			name: "error_no_candidate",
			input: PullRequestReassignReq{
				PullRequestId: "pr-2001",
				OldReviewerId: "u3",
			},
			userRepoSetup: func(repo *repoMocks.MockUserRepository) {
				repo.EXPECT().GetById(gomock.Any(), "u3").Return(domain.User{
					Id:       "u3",
					Name:     "test",
					IsActive: true,
					TeamId:   1,
				}, nil)

				repo.EXPECT().
					GetReassignCandidates(gomock.Any(), "u1", gomock.Any(), gomock.Any()).
					Return([]domain.User{}, nil)
			},
			prRepoSetup: func(repo *repoMocks.MockPullRequestRepository) {
				repo.EXPECT().
					GetByPrIdWithReviewersIds(gomock.Any(), "pr-2001").
					Return(r.GetByPrIdWithReviewersIdsDTO{
						Pr: domain.PullRequest{
							Id:                "pr-2001",
							Name:              "My PR",
							AuthorId:          "u1",
							StatusId:          1,
							NeedMoreReviewers: false,
							CreatedAt:         time.Now(),
							MergedAt:          nil,
						},
						ReviewersIds: []string{"u2", "u3"},
						StatusName:   domain.OPEN,
					}, nil)
			},
			reviewerRepoSetup: func(repo *repoMocks.MockPrReviewerRepository) {},
			expectedRes:       PullRequestReassignRes{},
			expectedErr:       e.ErrPrNoCandidate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userRepoMock := repoMocks.NewMockUserRepository(ctrl)
			prRepoMock := repoMocks.NewMockPullRequestRepository(ctrl)
			reviewerRepoMock := repoMocks.NewMockPrReviewerRepository(ctrl)

			tt.userRepoSetup(userRepoMock)
			tt.prRepoSetup(prRepoMock)
			tt.reviewerRepoSetup(reviewerRepoMock)

			uc := PullRequestUseCase{
				userRepo:     userRepoMock,
				prRepo:       prRepoMock,
				reviewerRepo: reviewerRepoMock,
			}

			res, err := uc.ReviewerReassign(context.Background(), tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error: got %v, want %v", err, tt.expectedErr)
			}

			res.Pr.CreatedAt = nil
			res.Pr.MergedAt = nil
			tt.expectedRes.Pr.CreatedAt = nil
			tt.expectedRes.Pr.MergedAt = nil

			require.Equal(t, tt.expectedRes, res)

		})
	}
}

func ptr(s time.Time) *time.Time {
	return &s
}
