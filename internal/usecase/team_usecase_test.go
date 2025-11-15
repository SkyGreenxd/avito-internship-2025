package usecase

import (
	"avito-internship/internal/domain"
	"avito-internship/pkg/e"
	"context"
	"errors"
	"testing"

	repoMocks "avito-internship/internal/repository/mocks"
	trMock "avito-internship/pkg/transaction/mocks"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTeamUseCase_AddTeam(t *testing.T) {
	tests := []struct {
		name          string
		input         TeamAddReq
		teamRepoSetup func(*repoMocks.MockTeamRepository)
		userRepoSetup func(*repoMocks.MockUserRepository)
		expectedRes   TeamAddRes
		expectedErr   error
	}{
		{
			name: "success",
			input: TeamAddReq{
				TeamName: "test",
				Members: []TeamMemberDTO{
					{
						Id:       "u1",
						Username: "test1",
						IsActive: true,
					},
				},
			},
			teamRepoSetup: func(teamRepo *repoMocks.MockTeamRepository) {
				teamRepo.EXPECT().
					Create(
						gomock.Any(),
						domain.NewTeam("test"),
					).
					Return(domain.Team{
						1,
						"test",
					}, nil)
			},
			userRepoSetup: func(userRepo *repoMocks.MockUserRepository) {
				userRepo.EXPECT().AddUsersToTeam(
					gomock.Any(),
					1,
					[]domain.User{
						{
							Id:       "u1",
							Name:     "test1",
							IsActive: true,
						},
					}).
					Return([]domain.User{
						{
							Id:       "u1",
							Name:     "test1",
							IsActive: true,
						},
					}, nil)
			},

			expectedRes: TeamAddRes{
				Team: TeamDTO{
					TeamName: "test",
					Members: []TeamMemberDTO{
						{
							Id:       "u1",
							Username: "test1",
							IsActive: true,
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name: "team exists",
			input: TeamAddReq{
				TeamName: "test",
				Members: []TeamMemberDTO{
					{
						Id:       "u1",
						Username: "test1",
						IsActive: true,
					},
				},
			},
			teamRepoSetup: func(teamRepo *repoMocks.MockTeamRepository) {
				teamRepo.EXPECT().
					Create(gomock.Any(), domain.NewTeam("test")).
					Return(domain.Team{}, e.ErrTeamIsExists)
			},
			userRepoSetup: func(userRepo *repoMocks.MockUserRepository) {},
			expectedRes:   TeamAddRes{},
			expectedErr:   e.ErrTeamIsExists,
		},
		{
			name: "member list is empty",
			input: TeamAddReq{
				TeamName: "test",
				Members:  []TeamMemberDTO{},
			},
			teamRepoSetup: func(teamRepo *repoMocks.MockTeamRepository) {},
			userRepoSetup: func(userRepo *repoMocks.MockUserRepository) {},
			expectedRes:   TeamAddRes{},
			expectedErr:   e.ErrEmptyMembers,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			teamRepo := repoMocks.NewMockTeamRepository(ctrl)
			userRepo := repoMocks.NewMockUserRepository(ctrl)

			mockTx := trMock.NewMockTx(ctrl)
			mockTx.EXPECT().Commit(gomock.Any()).Return(nil).AnyTimes()
			mockTx.EXPECT().Rollback(gomock.Any()).Return(nil).AnyTimes()

			mockTxPool := trMock.NewMockTransactional(ctrl)
			mockTxPool.EXPECT().
				BeginTx(gomock.Any(), gomock.Any()).
				Return(mockTx, nil).
				AnyTimes()

			teamUC := NewTeamUseCase(teamRepo, userRepo, mockTxPool)
			tt.teamRepoSetup(teamRepo)
			tt.userRepoSetup(userRepo)

			res, err := teamUC.AddTeam(context.Background(), tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error: got %v, want %v", err, tt.expectedErr)
			}

			require.Equal(t, tt.expectedRes, res)
		})
	}
}

func TestTeamUseCase_GetTeam(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		teamRepoSetup func(*repoMocks.MockTeamRepository)
		expectedRes   GetTeamRes
		expectedErr   error
	}{
		{
			name:  "success",
			input: "Team Name",
			teamRepoSetup: func(teamRepo *repoMocks.MockTeamRepository) {
				teamRepo.EXPECT().GetMembersByTeamNameWithUsers(
					gomock.Any(),
					"Team Name",
				).Return([]domain.User{
					{
						Id:       "u1",
						Name:     "test1",
						IsActive: true,
					},
					{
						Id:       "u2",
						Name:     "test2",
						IsActive: true,
					},
				}, nil)
			},
			expectedRes: GetTeamRes{
				TeamName: "Team Name",
				Members: []TeamMemberDTO{
					{
						Id:       "u1",
						Username: "test1",
						IsActive: true,
					},
					{
						Id:       "u2",
						Username: "test2",
						IsActive: true,
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:  "team not found",
			input: "Team Name888",
			teamRepoSetup: func(teamRepo *repoMocks.MockTeamRepository) {
				teamRepo.EXPECT().GetMembersByTeamNameWithUsers(
					gomock.Any(),
					"Team Name888",
				).Return([]domain.User{}, e.ErrTeamNotFound)
			},
			expectedRes: GetTeamRes{},
			expectedErr: e.ErrTeamNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			teamRepo := repoMocks.NewMockTeamRepository(ctrl)
			userRepo := repoMocks.NewMockUserRepository(ctrl)

			mockTx := trMock.NewMockTx(ctrl)
			mockTx.EXPECT().Commit(gomock.Any()).Return(nil).AnyTimes()
			mockTx.EXPECT().Rollback(gomock.Any()).Return(nil).AnyTimes()

			mockTxPool := trMock.NewMockTransactional(ctrl)
			mockTxPool.EXPECT().
				BeginTx(gomock.Any(), gomock.Any()).
				Return(mockTx, nil).
				AnyTimes()

			teamUC := NewTeamUseCase(teamRepo, userRepo, mockTxPool)
			tt.teamRepoSetup(teamRepo)

			res, err := teamUC.GetTeam(context.Background(), tt.input)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("unexpected error: got %v, want %v", err, tt.expectedErr)
			}

			require.Equal(t, tt.expectedRes, res)
		})
	}
}
