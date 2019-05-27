package db

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/hideUW/nuxt_go_template/server/domain/model"
	"github.com/hideUW/nuxt_go_template/server/domain/repository"
	"github.com/hideUW/nuxt_go_template/server/testutil"
	"github.com/pkg/errors"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	sessionValidIDForTest   = "testValidSessionID12345678"
	sessionInValidIDForTest = "testInvalidSessionID12345678"
)

func TestNewSessionRepository(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name string
		args args
		want repository.SessionRepository
	}{
		{
			name: "When it is given valid args, returns SessionRepository.",
			args: args{
				ctx: context.Background(),
			},
			want: &sessionRepository{
				context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSessionRepository(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSessionRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sessionRepository_ErroMsg(t *testing.T) {
	type fields struct {
		ctx context.Context
	}

	type args struct {
		method model.RepositoryMethod
		err    error
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr *model.RepositoryError
	}{
		{
			name: "When it is given valid args, returns appropriate error.",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				method: model.RepositoryMethodINSERT,
				err:    errors.New(errMsg),
			},
			wantErr: &model.RepositoryError{
				BaseErr:                     errors.New(errMsg),
				RepositoryMethod:            model.RepositoryMethodINSERT,
				DomainModelNameForDeveloper: model.DomainModelNameSessionForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameSessionForUser,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &sessionRepository{
				ctx: tt.fields.ctx,
			}
			if err := repo.ErrorMsg(tt.args.method, tt.args.err); errors.Cause(err).Error() != tt.wantErr.Error() {
				t.Errorf("sessionRepository.ErrMsg() error = %#v, want %#v", err, tt.wantErr)
			}
		})
	}
}

func Test_sessionRepository_GetSessionByID(t *testing.T) {
	// Set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		return
	}

	defer db.Close()

	testutil.SetFakeTime(time.Now())

	type fields struct {
		ctx context.Context
	}

	type args struct {
		m  repository.DBManager
		id string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Session
		wantErr *model.NoSuchDataError
	}{
		{
			name: "When a session specified by id exists, return a session.",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m:  db,
				id: sessionValidIDForTest,
			},
			want: &model.Session{
				ID:        sessionValidIDForTest,
				UserID:    userValidIDForTest,
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},
			wantErr: nil,
		},
		{
			name: "When a session specified by id does not exist, returns NoSuchDataError",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m:  db,
				id: sessionInValidIDForTest,
			},
			want: nil,
			wantErr: &model.NoSuchDataError{
				PropertyNameForDeveloper:    model.IDPropertyForDeveloper,
				PropertyNameForUser:         model.IDPropertyForUser,
				PropertyValue:               sessionInValidIDForTest,
				DomainModelNameForDeveloper: model.DomainModelNameSessionForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameSessionForUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := "SELECT id, user_id, created_at, updated_at FROM sessions WHERE id=?"
			prep := mock.ExpectPrepare(q)

			if tt.wantErr != nil {
				prep.ExpectQuery().WillReturnError(tt.wantErr)
			} else {
				rows := sqlmock.NewRows([]string{"id", "user_id", "created_at", "updated_at"}).
					AddRow(tt.want.ID, tt.want.UserID, tt.want.CreatedAt, tt.want.UpdatedAt)
				prep.ExpectQuery().WithArgs(tt.want.ID).WillReturnRows(rows)
			}

			repo := &sessionRepository{
				ctx: tt.fields.ctx,
			}
			got, err := repo.GetSessionByID(tt.args.m, tt.args.id)

			if tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("sessionRepository.GetSessionByID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sessionRepository.GetSessionByID() = %v, wantErr %v", got, tt.want)
			}
		})
	}
}

func Test_sessionRepository_InsertSession(t *testing.T) {
	// set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	defer db.Close()

	testutil.SetFakeTime(time.Now())

	type fields struct {
		ctx context.Context
	}
	type args struct {
		m       repository.DBManager
		session *model.Session
		err     error
	}

	tests := []struct {
		name        string
		fields      fields
		args        args
		rowAffected int64
		want        string
		wantErr     *model.RepositoryError
	}{
		{
			name: "When a session which has ID, User_ID, CreatedAt, UpdatedAt is given, returns ID",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m: db,
				session: &model.Session{
					ID:        sessionValidIDForTest,
					UserID:    userValidIDForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 1,
			wantErr:     nil,
		},
		{
			name: "when RowAffected is 0、returns error",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m: db,
				session: &model.Session{
					ID:        sessionInValidIDForTest,
					UserID:    userValidIDForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodINSERT,
				DomainModelNameForDeveloper: model.DomainModelNameSessionForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameSessionForUser,
			},
		},
		{
			name: "when RowAffected is 2、returns error",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m: db,
				session: &model.Session{
					ID:        sessionInValidIDForTest,
					UserID:    userValidIDForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 2,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodINSERT,
				DomainModelNameForDeveloper: model.DomainModelNameSessionForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameSessionForUser,
			},
		},
		{
			name: "when DB error has occurred、returns error",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m: db,
				session: &model.Session{
					ID:        sessionInValidIDForTest,
					UserID:    userValidIDForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: errors.New(errMsg),
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodINSERT,
				DomainModelNameForDeveloper: model.DomainModelNameSessionForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameSessionForUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "INSERT INTO sessions"
			prep := mock.ExpectPrepare(query)

			if tt.args.err != nil {
				prep.ExpectExec().WithArgs(tt.args.session.ID, tt.args.session.UserID, tt.args.session.CreatedAt, tt.args.session.UpdatedAt).WillReturnError(tt.args.err)
			} else {
				prep.ExpectExec().WithArgs(tt.args.session.ID, tt.args.session.UserID, tt.args.session.CreatedAt, tt.args.session.UpdatedAt).WillReturnResult(sqlmock.NewResult(1, tt.rowAffected))
			}

			repo := &sessionRepository{
				ctx: tt.fields.ctx,
			}

			err := repo.InsertSession(tt.args.m, tt.args.session)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("sessionRepository.InsertSession() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}
