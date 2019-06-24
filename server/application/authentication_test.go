package application

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	"github.com/hideUW/nuxt-go-chat-app/server/domain/model"
	"github.com/hideUW/nuxt-go-chat-app/server/domain/repository"
	mock_repository "github.com/hideUW/nuxt-go-chat-app/server/domain/repository/mock"
	"github.com/hideUW/nuxt-go-chat-app/server/domain/service"
	mock_service "github.com/hideUW/nuxt-go-chat-app/server/domain/service/mock"
	"github.com/hideUW/nuxt-go-chat-app/server/testutil"
)

func Test_authenticationService_SignUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		m                 repository.DBManager
		userRepository    repository.UserRepository
		sessionRepository repository.SessionRepository
		userService       service.UserService
		sessionService    service.SessionService
		txCloser          CloseTransaction
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}

	type mockUserRepoArgs struct {
		user *model.User
	}

	type mockSessionRepoArgs struct {
		session *model.Session
	}

	type mockUserServiceArgs struct {
		ctx  context.Context
		name string
	}

	type mockSessionServiceArgs struct {
		ctx context.Context
		id  string
	}

	type mockUserRepoReturns struct {
		id  uint32
		err error
	}

	type mockUserServiceReturns struct {
		found bool
		err   error
	}

	type mockSessionRepoReturns struct {
		err error
	}

	type mockSessionServiceReturns struct {
		found bool
		err   error
	}

	user, err := model.NewUser(model.UserNameForTest, model.PasswordForTest)
	if err != nil {
		t.Fatal(err)
	}

	userForSession, err := model.NewUser(model.UserNameForTest, model.PasswordForTest)
	if err != nil {
		t.Fatal(err)
	}
	userForSession.Password = user.Password
	userForSession.SessionID = model.SessionValidIDForTest

	tests := []struct {
		name   string
		fields fields
		args   args
		mockUserRepoArgs
		mockUserServiceArgs
		mockSessionRepoArgs
		mockSessionServiceArgs
		mockUserRepoReturns
		mockUserServiceReturns
		mockSessionRepoReturns
		mockSessionServiceReturns
		wantUser *model.User
		wantErr  error
	}{
		// TODO: Add test cases.
		{
			name: "When appropriate name and password are given and the user which name is same as given name does'nt exist, returns user and nil",
			fields: fields{
				m:                 mock_repository.NewMockDBManager(ctrl),
				userRepository:    mock_repository.NewMockUserRepository(ctrl),
				sessionRepository: mock_repository.NewMockSessionRepository(ctrl),
				userService:       mock_service.NewMockUserService(ctrl),
				sessionService:    mock_service.NewMockSessionService(ctrl),
				txCloser: func(tx repository.TxManager, err error) error {
					return nil
				},
			},
			args: args{
				ctx:  context.Background(),
				user: user,
			},
			mockUserRepoArgs: mockUserRepoArgs{
				user: userForSession,
			},
			mockSessionRepoArgs: mockSessionRepoArgs{
				session: &model.Session{
					ID:        model.SessionValidIDForTest,
					UserID:    model.UserValidIDForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			mockUserServiceArgs: mockUserServiceArgs{
				ctx:  context.Background(),
				name: model.UserNameForTest,
			},
			mockSessionServiceArgs: mockSessionServiceArgs{
				ctx: context.Background(),
				id:  model.SessionValidIDForTest,
			},
			mockUserRepoReturns: mockUserRepoReturns{
				id:  model.UserValidIDForTest,
				err: nil,
			},
			mockSessionRepoReturns: mockSessionRepoReturns{
				err: nil,
			},
			mockUserServiceReturns: mockUserServiceReturns{
				found: false,
				err:   nil,
			},
			mockSessionServiceReturns: mockSessionServiceReturns{
				found: false,
				err:   nil,
			},
			wantUser: &model.User{
				ID:        model.UserValidIDForTest,
				Name:      model.UserNameForTest,
				Password:  model.PasswordForTest,
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},

			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, ok := tt.fields.m.(*mock_repository.MockDBManager)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}
			m.EXPECT().Begin().Return(mock_repository.NewMockTxManager(ctrl), nil)

			us, ok := tt.fields.userService.(*mock_service.MockUserService)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}
			us.EXPECT().IsAlreadyExistName(tt.args.ctx, tt.mockUserServiceArgs.name).Return(tt.mockUserServiceReturns.found, tt.mockUserServiceReturns.err)

			ur, ok := tt.fields.userRepository.(*mock_repository.MockUserRepository)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}
			ur.EXPECT().InsertUser(tt.fields.m, tt.mockUserRepoArgs.user).Return(tt.mockUserRepoReturns.id, tt.mockUserRepoReturns.err)

			ss, ok := tt.fields.sessionService.(*mock_service.MockSessionService)
			if !ok {
				t.Fatal("failed to assert MockUserRepository")
			}
			ss.EXPECT().IsAlreadyExistID(tt.mockSessionServiceArgs.ctx, tt.mockSessionServiceArgs.id).Return(tt.mockSessionServiceReturns.found, tt.mockSessionServiceReturns.err)

			ss.EXPECT().SessionID().Return(model.SessionValidIDForTest)

			sr, ok := tt.fields.sessionRepository.(*mock_repository.MockSessionRepository)
			if !ok {
				t.Fatal("failed to assert MockSessionRepository")
			}
			sr.EXPECT().InsertSession(tt.fields.m, tt.mockSessionRepoArgs.session).Return(tt.mockSessionRepoReturns.err)

			a := &authenticationService{
				m:                 tt.fields.m,
				userRepository:    tt.fields.userRepository,
				sessionRepository: tt.fields.sessionRepository,
				userService:       tt.fields.userService,
				sessionService:    tt.fields.sessionService,
				txCloser:          tt.fields.txCloser,
			}
			gotUser, err := a.SignUp(tt.args.ctx, tt.args.user)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("authenticationService.SignUp() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("authenticationService.SignUp() = %v, want %v", gotUser, tt.wantUser)
			}
		})
	}
}
