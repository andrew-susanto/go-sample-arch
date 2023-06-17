package psql

import (
	// golang package
	"context"
	"testing"

	// external package
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRepository_GetUserByID(t *testing.T) {
	type args struct {
		userID int64
	}

	tests := []struct {
		name     string
		args     args
		want     User
		wantErr  error
		mockFunc func(mock sqlmock.Sqlmock)
	}{
		{
			name: "when_get_user_by_id_success_then_xxx",
			args: args{
				userID: 123,
			},
			want: User{
				ID: 123,
			},
			wantErr: nil,
			mockFunc: func(mock sqlmock.Sqlmock) {
				rows := mock.NewRows([]string{"Id"}).
					AddRow(123)

				mock.ExpectQuery("SELECT Id FROM Users WHERE Id = ?").WithArgs(123).WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db, mock, _ := sqlmock.New()
			defer db.Close()

			repository := NewRepository(db)

			tt.mockFunc(mock)
			got, gotErr := repository.GetUserByID(context.Background(), tt.args.userID)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)

			// we make sure that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
