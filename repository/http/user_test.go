package http

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
		mockFunc func(mock *MockHttpClient)
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
			mockFunc: func(mock *MockHttpClient) {

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHttpClient := NewMockHttpClient(ctrl)
			repository := NewRepository(mockHttpClient)

			tt.mockFunc(mockHttpClient)
			got, gotErr := repository.GetUserByID(context.Background(), tt.args.userID)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
