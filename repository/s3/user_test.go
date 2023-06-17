package s3

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
		mockFunc func(s3Service *MockS3Service)
	}{
		{
			name: "when_get_user_by_id_success_then_xxx",
			args: args{
				userID: 123,
			},
			want: User{
				ID:     123,
				Name:   "dummy user name",
				Gender: 3,
			},
			wantErr: nil,
			mockFunc: func(s3Service *MockS3Service) {
				s3Service.EXPECT()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockS3 := NewMockS3Service(ctrl)
			repository := NewRepository(mockS3)

			tt.mockFunc(mockS3)
			got, gotErr := repository.GetUserByID(context.Background(), tt.args.userID)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
