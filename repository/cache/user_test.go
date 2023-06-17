package cache

import (
	// golang package
	"context"
	"testing"

	// external package
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
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
		mockFunc func(mock *MockRedis)
	}{
		{
			name: "when_get_user_by_id_success_then_xxx",
			args: args{
				userID: 123,
			},
			want:    User{},
			wantErr: nil,
			mockFunc: func(mock *MockRedis) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRedis := NewMockRedis(ctrl)
			repository := NewRepository(mockRedis)

			tt.mockFunc(mockRedis)
			got, gotErr := repository.GetUserByID(context.Background(), tt.args.userID)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
