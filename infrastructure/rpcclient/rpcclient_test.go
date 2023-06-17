package rpcclient

import (
	// golang package
	"context"
	"errors"
	"testing"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/config"

	// external package
	"github.com/stretchr/testify/assert"
	"github.com/ybbus/jsonrpc/v3"
	"go.uber.org/mock/gomock"
)

func TestRpcClient_NewClient(t *testing.T) {
	type args struct {
		config config.RpcClientConfig
	}

	tests := []struct {
		name string
		args args
		want RpcClient
	}{
		{
			name: "when_given_config_then_return_rpcclient",
			args: args{
				config: config.RpcClientConfig{
					ServiceUrl:               "test",
					ServiceName:              "testing",
					Timeout:                  30,
					MaxConcurrentRequest:     10,
					ErrorPercentageThreshold: 25,
				},
			},
			want: RpcClient{
				cbConfig: map[string]bool{},
				client:   jsonrpc.NewClient("test"),
				config: config.RpcClientConfig{
					ServiceUrl:               "test",
					ServiceName:              "testing",
					Timeout:                  30,
					MaxConcurrentRequest:     10,
					ErrorPercentageThreshold: 25,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			got := NewClient(tt.args.config)
			assert.Equal(t, tt.want.config, got.config)
		})
	}
}

func TestRpcClient_Call(t *testing.T) {
	type args struct {
		method string
		params interface{}
		config config.RpcClientConfig
	}

	tests := []struct {
		name    string
		args    args
		want    jsonrpc.RPCResponse
		mock    func(client *MockjsonRpcClient)
		wantErr error
	}{
		{
			name: "when_given_client_name_not_exists_and_return_success_expect_return_success",
			args: args{
				config: config.RpcClientConfig{
					ServiceUrl:               "test",
					ServiceName:              "testing",
					Timeout:                  30,
					MaxConcurrentRequest:     10,
					ErrorPercentageThreshold: 25,
				},
				method: "method_test",
				params: "abc",
			},
			mock: func(client *MockjsonRpcClient) {
				client.EXPECT().Call(gomock.Any(), "method_test", "abc").Return(&jsonrpc.RPCResponse{
					Result: "test",
				}, nil)
			},
			want: jsonrpc.RPCResponse{
				Result: "test",
			},
			wantErr: nil,
		},
		{
			name: "when_given_client_name_not_exists_and_return_error_expect_return_error",
			args: args{
				config: config.RpcClientConfig{
					ServiceUrl:               "test",
					ServiceName:              "testing",
					Timeout:                  30,
					MaxConcurrentRequest:     10,
					ErrorPercentageThreshold: 25,
				},
				method: "method_test",
				params: "abc",
			},
			mock: func(client *MockjsonRpcClient) {
				client.EXPECT().Call(gomock.Any(), "method_test", "abc").Return(&jsonrpc.RPCResponse{
					Result: "test",
				}, errors.New("some_error"))
			},
			want: jsonrpc.RPCResponse{
				Result: "test",
			},
			wantErr: errors.New("some_error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			jsonRpcClient := NewClient(tt.args.config)

			jsonRpcClientMock := NewMockjsonRpcClient(ctrl)
			jsonRpcClient.client = jsonRpcClientMock

			tt.mock(jsonRpcClientMock)
			got, gotErr := jsonRpcClient.Call(context.Background(), tt.args.method, tt.args.params)

			assert.Equal(t, tt.want, *got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
