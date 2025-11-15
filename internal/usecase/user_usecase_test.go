package usecase

import (
	"avito-internship/internal/domain"
	r "avito-internship/internal/repository"
	"avito-internship/internal/repository/mocks"
	"avito-internship/pkg/e"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserUseCase_SetIsActive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reviewerRepo := mocks.NewMockPrReviewerRepository(ctrl)
	userRepo := mocks.NewMockUserRepository(ctrl)
	teamRepo := mocks.NewMockTeamRepository(ctrl)

	userUC := NewUserUseCase(reviewerRepo, userRepo, teamRepo)

	tests := []struct {
		name          string
		input         SetIsActiveReq
		userRepoSetup func(*mocks.MockUserRepository)
		teamRepoSetup func(*mocks.MockTeamRepository)
		expectedRes   SetIsActiveRes
		expectedErr   error
	}{
		{
			name: "success",
			input: SetIsActiveReq{
				UserId:   "u2",
				IsActive: false,
			},
			userRepoSetup: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().
					UpdateIsActive(gomock.Any(), "u2", false).
					Return(domain.User{
						Id:       "u2",
						Name:     "Test User",
						IsActive: false,
						TeamId:   1,
					}, nil)
			},
			teamRepoSetup: func(teamRepo *mocks.MockTeamRepository) {
				teamRepo.EXPECT().
					GetTeamByUserId(gomock.Any(), "u2").
					Return(domain.Team{
						Id:   1,
						Name: "Test Team",
					}, nil)
			},
			expectedRes: SetIsActiveRes{
				User: UserDTO{
					Id:       "u2",
					Username: "Test User",
					IsActive: false,
					TeamName: "Test Team",
				},
			},
			expectedErr: nil,
		},
		{
			name: "user not found",
			input: SetIsActiveReq{
				UserId:   "u888",
				IsActive: false,
			},
			userRepoSetup: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().
					UpdateIsActive(gomock.Any(), "u888", false).
					Return(domain.User{}, e.ErrUserNotFound)
			},
			teamRepoSetup: func(teamRepo *mocks.MockTeamRepository) {
			},
			expectedRes: SetIsActiveRes{},
			expectedErr: e.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.userRepoSetup(userRepo)
			tt.teamRepoSetup(teamRepo)

			res, err := userUC.SetIsActive(context.Background(), tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error: got %v, want %v", err, tt.expectedErr)
			}

			require.Equal(t, tt.expectedRes, res)
		})
	}
}

func TestUserUseCase_GetReview(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		userRepoSetup     func(*mocks.MockUserRepository)
		reviewerRepoSetup func(*mocks.MockPrReviewerRepository)
		expectedRes       GetReviewRes
		expectedErr       error
	}{
		{
			name:  "success",
			input: "u1",
			userRepoSetup: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().GetById(gomock.Any(), "u1").
					Return(domain.User{
						Id:       "u1",
						Name:     "Test User",
						IsActive: true,
						TeamId:   1,
					}, nil)
			},
			reviewerRepoSetup: func(userRepo *mocks.MockPrReviewerRepository) {
				userRepo.EXPECT().GetPRByReviewer(gomock.Any(), "u1").
					Return(r.GetPRByReviewerDTO{
						Prs: make([]r.PrWithStatusName, 0),
					}, nil)
			},
			expectedRes: GetReviewRes{
				UserId:       "u1",
				PullRequests: make([]PullRequestShort, 0),
			},
			expectedErr: nil,
		},
		{
			name:  "user not found",
			input: "u888",
			userRepoSetup: func(userRepo *mocks.MockUserRepository) {
				userRepo.EXPECT().GetById(gomock.Any(), "u888").
					Return(domain.User{}, e.ErrUserNotFound)
			},
			reviewerRepoSetup: func(userRepo *mocks.MockPrReviewerRepository) {},
			expectedRes:       GetReviewRes{},
			expectedErr:       e.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			reviewerRepo := mocks.NewMockPrReviewerRepository(ctrl)
			userRepo := mocks.NewMockUserRepository(ctrl)
			teamRepo := mocks.NewMockTeamRepository(ctrl)

			userUC := NewUserUseCase(reviewerRepo, userRepo, teamRepo)

			tt.userRepoSetup(userRepo)
			tt.reviewerRepoSetup(reviewerRepo)

			res, err := userUC.GetReview(context.Background(), tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error: got %v, want %v", err, tt.expectedErr)
			}

			require.Equal(t, tt.expectedRes, res)
		})
	}
}
