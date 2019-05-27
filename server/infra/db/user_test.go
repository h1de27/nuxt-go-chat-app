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
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	userNameForTest             = "testUserName"
	passwordForTest             = "testPasswor"
	userValidIDForTest   uint32 = 1
	userInValidIDForTest uint32 = 2
)

func TestNewUserRepository(t *testing.T) {
	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name string
		args args
		want repository.UserRepository
	}{
		{
			name: "When it is given appropriate args, return UserRepository.",
			args: args{
				ctx: context.Background(),
			},
			want: &userRepository{
				context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserRepository(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepository_ErrorMsg(t *testing.T) {

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
			name: "When it is given the appropriate args, returns appropriate error.",
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
				DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameUserForUser,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &userRepository{
				ctx: tt.fields.ctx,
			}
			if err := repo.ErrorMsg(tt.args.method, tt.args.err); errors.Cause(err).Error() != tt.wantErr.Error() {
				t.Errorf("userRepository.ErrorMsg() error = %#v, wantErr %#v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepository_GetUserByID(t *testing.T) {
	// Set sql mock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		return
	}

	// If db has error, close db.
	defer db.Close()

	// Set fake time.
	testutil.SetFakeTime(time.Now())

	type fields struct {
		ctx context.Context
	}

	type args struct {
		m  repository.DBManager
		id uint32
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr *model.NoSuchDataError
	}{
		{
			name: "When the specific user exists, return the user.",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m:  db,
				id: userValidIDForTest,
			},
			want: &model.User{
				ID:        userValidIDForTest,
				Name:      userNameForTest,
				SessionID: sessionValidIDForTest,
				Password:  passwordForTest,
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},
			wantErr: nil,
		},
		{
			name: "When the specific user does not exist, return NoSuchDataError",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m:  db,
				id: userInValidIDForTest,
			},
			want: nil,
			wantErr: &model.NoSuchDataError{
				PropertyNameForDeveloper:    model.IDPropertyForDeveloper,
				PropertyNameForUser:         model.IDPropertyForUser,
				PropertyValue:               userInValidIDForTest,
				DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameUserForUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := "SELECT id, name, session_id, password, created_at, updated_at FROM users WHERE id=?"
			prep := mock.ExpectPrepare(q)

			if tt.wantErr != nil {
				prep.ExpectQuery().WillReturnError(tt.wantErr)
			} else {
				rows := sqlmock.NewRows([]string{"id", "name", "session_id", "password", "created_at", "updated_at"}).
					AddRow(tt.want.ID, tt.want.Name, tt.want.SessionID, tt.want.Password, tt.want.CreatedAt, tt.want.UpdatedAt)
				prep.ExpectQuery().WithArgs(tt.want.ID).WillReturnRows(rows)
			}

			repo := &userRepository{
				ctx: tt.fields.ctx,
			}
			got, err := repo.GetUserByID(tt.args.m, tt.args.id)

			if tt.wantErr != nil {
				if !reflect.DeepEqual(err, tt.wantErr) {
					t.Errorf("userRepository.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepository.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepository_GetUserByName(t *testing.T) {
	// Set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		return
	}

	defer db.Close()

	type fields struct {
		ctx context.Context
	}

	type args struct {
		m    repository.DBManager
		name string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr *model.NoSuchDataError
	}{
		{
			name: "Return a user if a user specified by name exists.",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m:    db,
				name: userNameForTest,
			},
			want: &model.User{
				Name:      userNameForTest,
				SessionID: sessionValidIDForTest,
				Password:  passwordForTest,
				CreatedAt: testutil.TimeNow(),
				UpdatedAt: testutil.TimeNow(),
			},
			wantErr: nil,
		},
		{
			name: "Return NoSuchDataError if a user specified by name doesn't exists.",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m:    db,
				name: "test2",
			},
			want: nil,
			wantErr: &model.NoSuchDataError{
				PropertyNameForDeveloper:    model.NamePropertyForDeveloper,
				PropertyNameForUser:         model.NamePropertyForUser,
				PropertyValue:               "test2",
				DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameUserForUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := "SELECT id, name, session_id, password, created_at, updated_at FROM users WHERE name=?"
			prep := mock.ExpectPrepare(q)

			if tt.wantErr != nil {
				prep.ExpectQuery().WillReturnError(tt.wantErr)
			} else {
				rows := sqlmock.NewRows([]string{"id", "name", "session_id", "password", "created_at", "updated_at"}).
					AddRow(tt.want.ID, tt.want.Name, tt.want.SessionID, tt.want.Password, tt.want.CreatedAt, tt.want.UpdatedAt)
				prep.ExpectQuery().WithArgs(tt.want.Name).WillReturnRows(rows)
			}

			repo := &userRepository{
				ctx: tt.fields.ctx,
			}
			got, err := repo.GetUserByName(tt.args.m, tt.args.name)
			if tt.wantErr != nil {
				if !reflect.DeepEqual(err, tt.wantErr) {
					t.Errorf("userRepository.GetUserByName() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepository.GetUserByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepository_InsertUser(t *testing.T) {
	// Set sql mock
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
		m    repository.DBManager
		user *model.User
		err  error
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		rowAffected int64
		wantErr     *model.RepositoryError
	}{
		{
			name: "When a user which has sID, Name, Session_id, Password, CreatedAt, UpdatedAt is given, returns nil",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m: db,
				user: &model.User{
					ID:        userValidIDForTest,
					Name:      userNameForTest,
					SessionID: sessionValidIDForTest,
					Password:  passwordForTest,
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
				user: &model.User{
					ID:        userInValidIDForTest,
					Name:      userNameForTest,
					SessionID: sessionValidIDForTest,
					Password:  passwordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodINSERT,
				DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameUserForUser,
			},
		},
		{
			name: "when RowAffected is 2、returns error",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m: db,
				user: &model.User{
					ID:        userInValidIDForTest,
					Name:      userNameForTest,
					SessionID: sessionValidIDForTest,
					Password:  passwordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 2,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodINSERT,
				DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameUserForUser,
			},
		},
		{
			name: "when DB error has occurred、returns error",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m: db,
				user: &model.User{
					ID:        userInValidIDForTest,
					Name:      userNameForTest,
					SessionID: sessionValidIDForTest,
					Password:  passwordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: errors.New(errMsg),
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodINSERT,
				DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameUserForUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "INSERT INTO users"
			prep := mock.ExpectPrepare(query)

			if tt.args.err != nil {
				prep.ExpectExec().WithArgs(tt.args.user.ID, tt.args.user.Name, tt.args.user.SessionID, tt.args.user.Password, tt.args.user.CreatedAt, tt.args.user.UpdatedAt).WillReturnError(tt.args.err)
			} else {
				prep.ExpectExec().WithArgs(tt.args.user.ID, tt.args.user.Name, tt.args.user.SessionID, tt.args.user.Password, tt.args.user.CreatedAt, tt.args.user.UpdatedAt).WillReturnResult(sqlmock.NewResult(1, tt.rowAffected))
			}

			repo := &userRepository{
				ctx: tt.fields.ctx,
			}

			_, err := repo.InsertUser(tt.args.m, tt.args.user)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("userRepository.InsertUser() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}

}

func Test_userRepository_UpdateUser(t *testing.T) {
	// Set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	testutil.SetFakeTime(time.Now())

	type fields struct {
		ctx context.Context
	}

	type args struct {
		m    repository.DBManager
		id   uint32
		user *model.User
		err  error
	}

	tests := []struct {
		name        string
		fields      fields
		args        args
		rowAffected int64
		wantErr     *model.RepositoryError
	}{
		{
			name: "When a user which has Name, Session_id, Password, UpdatedAt is given, returns nil",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m:  db,
				id: userValidIDForTest,
				user: &model.User{
					ID:        userValidIDForTest,
					Name:      userNameForTest,
					SessionID: sessionValidIDForTest,
					Password:  passwordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: errors.New(errMsg),
			},
			rowAffected: 1,
			wantErr:     nil,
		},
		{
			name: "When RowAffected is 0, returns error",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m:  db,
				id: userInValidIDForTest,
				user: &model.User{
					ID:        userInValidIDForTest,
					Name:      userNameForTest,
					SessionID: sessionValidIDForTest,
					Password:  passwordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
				err: errors.New(errMsg),
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodUPDATE,
				DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameUserForUser,
			},
		},
		{
			name: "When RowAffected is 2, returns error.",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m:  db,
				id: userInValidIDForTest,
				user: &model.User{
					ID:        userInValidIDForTest,
					Name:      userNameForTest,
					SessionID: sessionValidIDForTest,
					Password:  passwordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 2,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodUPDATE,
				DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameUserForUser,
			},
		},
		{
			name: "when DB error has occurred、returns error",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				m:  db,
				id: userInValidIDForTest,
				user: &model.User{
					ID:        userInValidIDForTest,
					Name:      userNameForTest,
					SessionID: sessionValidIDForTest,
					Password:  passwordForTest,
					CreatedAt: testutil.TimeNow(),
					UpdatedAt: testutil.TimeNow(),
				},
			},
			rowAffected: 0,
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodUPDATE,
				DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameUserForUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "UPDATE users SET session_id=\\?, password=\\?, created_at=\\?, updated_at=\\? WHERE id=\\?"
			prep := mock.ExpectPrepare(query)

			if tt.args.err != nil {
				prep.ExpectExec().WithArgs(tt.args.user.SessionID, tt.args.user.Password, tt.args.user.CreatedAt, tt.args.user.UpdatedAt, tt.args.id).WillReturnError(tt.args.err)
			} else {
				prep.ExpectExec().WithArgs(tt.args.user.SessionID, tt.args.user.Password, tt.args.user.CreatedAt, tt.args.user.UpdatedAt, tt.args.id).WillReturnResult(sqlmock.NewResult(1, tt.rowAffected))
			}

			repo := &userRepository{
				ctx: tt.fields.ctx,
			}
			err := repo.UpdateUser(tt.args.m, tt.args.id, tt.args.user)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("userRepository.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}

func Test_userRepository_DeleteUser(t *testing.T) {
	// Set sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	testutil.SetFakeTime(time.Now())

	type fields struct {
		ctx context.Context
	}

	type args struct {
		m   repository.DBManager
		id  uint32
		err error
	}

	tests := []struct {
		name        string
		fields      fields
		rowAffected int64
		args        args
		wantErr     *model.RepositoryError
	}{
		{
			name: "When a user specified by id exists, returns nil",
			fields: fields{
				ctx: context.Background(),
			},
			rowAffected: 1,
			args: args{
				m:  db,
				id: userValidIDForTest,
			},
			wantErr: nil,
		},
		{
			name: "when RowAffected is 0、returns error",
			fields: fields{
				ctx: context.Background(),
			},
			rowAffected: 0,
			args: args{
				m:  db,
				id: userInValidIDForTest,
			},
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodDELETE,
				DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameUserForUser,
			},
		},
		{
			name: "when RowAffected is 2、returns error",
			fields: fields{
				ctx: context.Background(),
			},
			rowAffected: 2,
			args: args{
				m:  db,
				id: userInValidIDForTest,
			},
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodDELETE,
				DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameUserForUser,
			},
		},
		{
			name: "when DB error has occurred、returns error",
			fields: fields{
				ctx: context.Background(),
			},
			rowAffected: 0,
			args: args{
				m:   db,
				id:  userInValidIDForTest,
				err: errors.New(errMsg),
			},
			wantErr: &model.RepositoryError{
				RepositoryMethod:            model.RepositoryMethodDELETE,
				DomainModelNameForDeveloper: model.DomainModelNameUserForDeveloper,
				DomainModelNameForUser:      model.DomainModelNameUserForUser,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := "DELETE FROM users WHERE id=\\?"
			prep := mock.ExpectPrepare(query)

			if tt.args.err != nil {
				prep.ExpectExec().WithArgs(tt.args.id).WillReturnError(tt.args.err)
			} else {
				prep.ExpectExec().WithArgs(tt.args.id).WillReturnResult(sqlmock.NewResult(1, tt.rowAffected))
			}

			repo := &userRepository{
				ctx: tt.fields.ctx,
			}

			err := repo.DeleteUser(tt.args.m, tt.args.id)
			if tt.wantErr != nil {
				if errors.Cause(err).Error() != tt.wantErr.Error() {
					t.Errorf("userRepository.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}
