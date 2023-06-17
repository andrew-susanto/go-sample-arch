package jsonrpc

import (
	// golang package
	"context"
	"testing"

	// extenral package
	"github.com/stretchr/testify/assert"
	"github.com/ybbus/jsonrpc/v3"
	gomock "go.uber.org/mock/gomock"
)

func TestRepository_GetUserByID(t *testing.T) {
	type args struct {
		bookingID int64
	}

	tests := []struct {
		name     string
		args     args
		want     TripItinerary
		wantErr  error
		mockFunc func(cxpigw *MockRPCClient)
	}{
		{
			name: "when_get_user_by_id_success_then_xxx",
			args: args{
				bookingID: 123,
			},
			want:    TripItinerary{},
			wantErr: nil,
			mockFunc: func(cxpigw *MockRPCClient) {
				cxpigw.EXPECT().Call(gomock.Any(), "getTripItinerary", int64(123)).Return(&jsonrpc.RPCResponse{}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCxpIgwRpcClient := NewMockRPCClient(ctrl)
			repository := NewRepository(mockCxpIgwRpcClient)

			tt.mockFunc(mockCxpIgwRpcClient)
			got, gotErr := repository.GetTripItinerary(context.Background(), tt.args.bookingID)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
