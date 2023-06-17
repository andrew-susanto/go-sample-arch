package jsonrpccontext

import (
	// golang package
	"context"
	"net/http"
	"testing"

	// external package
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestJsonRpcContext_WriteHTTPResponseToJSON(t *testing.T) {
	type args struct {
		response             JSONRpcResponseSchema
		statusCode           int
		responseAlreadyWrite bool
	}

	tests := []struct {
		name string
		args args
		mock func(mockWriter *MockResponseWriter)
	}{
		{
			name: "when_given_json_rpc_response_then_expect_write_response",
			args: args{
				response: JSONRpcResponseSchema{
					JSONRpcVersion: "2.0",
					Result: map[string]interface{}{
						"testing": "testing",
					},
					Error: nil,
					ID:    123,
				},
				statusCode: 200,
			},
			mock: func(mockWriter *MockResponseWriter) {
				mockWriter.EXPECT().Header().Return(http.Header{})
				mockWriter.EXPECT().WriteHeader(200)
				mockWriter.EXPECT().Write([]byte("{\"jsonrpc\":\"2.0\",\"result\":{\"testing\":\"testing\"},\"id\":123}"))
			},
		},
		{
			name: "when_given_invalid_json_rpc_response_then_expect_write_error",
			args: args{
				response: JSONRpcResponseSchema{
					JSONRpcVersion: "2.0",
					Result: map[string]interface{}{
						"testing": "testing",
					},
					Error: nil,
					ID:    func() {},
				},
				statusCode: 200,
			},
			mock: func(mockWriter *MockResponseWriter) {
				mockWriter.EXPECT().Header().Return(http.Header{})
				mockWriter.EXPECT().WriteHeader(500)
				mockWriter.EXPECT().Write([]byte("{\"jsonrpc\":\"2.0\",\"error\":{\"code\":-32603,\"message\":\"Internal error\"},\"id\":null}"))
			},
		},
		{
			name: "when_respnse_already_write_then_expect_no_write_executed",
			args: args{
				response: JSONRpcResponseSchema{
					JSONRpcVersion: "2.0",
					Result: map[string]interface{}{
						"testing": "testing",
					},
					Error: nil,
					ID:    func() {},
				},
				statusCode:           200,
				responseAlreadyWrite: true,
			},
			mock: func(mockWriter *MockResponseWriter) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWriter := NewMockResponseWriter(ctrl)

			tdkHttpContext := TdkJsonRpcContext{
				Context:              context.Background(),
				Request:              nil,
				Writer:               mockWriter,
				ResponseAlreadyWrite: tt.args.responseAlreadyWrite,
			}
			tt.mock(mockWriter)

			tdkHttpContext.WriteHTTPResponseToJSON(tt.args.response, tt.args.statusCode)
		})
	}
}

func TestJsonRpcContext_WriteHTTPResponseBatchToJSON(t *testing.T) {
	type args struct {
		response             []JSONRpcResponseSchema
		statusCode           int
		responseAlreadyWrite bool
	}

	tests := []struct {
		name string
		args args
		mock func(mockWriter *MockResponseWriter)
	}{
		{
			name: "when_given_json_rpc_response_then_expect_write_response",
			args: args{
				response: []JSONRpcResponseSchema{{
					JSONRpcVersion: "2.0",
					Result: map[string]interface{}{
						"testing": "testing",
					},
					Error: nil,
					ID:    123,
				}},
				statusCode: 200,
			},
			mock: func(mockWriter *MockResponseWriter) {
				mockWriter.EXPECT().Header().Return(http.Header{})
				mockWriter.EXPECT().WriteHeader(200)
				mockWriter.EXPECT().Write([]byte("[{\"jsonrpc\":\"2.0\",\"result\":{\"testing\":\"testing\"},\"id\":123}]"))
			},
		},
		{
			name: "when_given_invalid_json_rpc_response_then_expect_write_error",
			args: args{
				response: []JSONRpcResponseSchema{{
					JSONRpcVersion: "2.0",
					Result: map[string]interface{}{
						"testing": "testing",
					},
					Error: nil,
					ID:    func() {},
				}},
				statusCode: 200,
			},
			mock: func(mockWriter *MockResponseWriter) {
				mockWriter.EXPECT().Header().Return(http.Header{})
				mockWriter.EXPECT().WriteHeader(500)
				mockWriter.EXPECT().Write([]byte("{\"jsonrpc\":\"2.0\",\"error\":{\"code\":-32603,\"message\":\"Internal error\"},\"id\":null}"))
			},
		},
		{
			name: "when_respnse_already_write_then_expect_no_write_executed",
			args: args{
				response: []JSONRpcResponseSchema{{
					JSONRpcVersion: "2.0",
					Result: map[string]interface{}{
						"testing": "testing",
					},
					Error: nil,
					ID:    func() {},
				}},
				statusCode:           200,
				responseAlreadyWrite: true,
			},
			mock: func(mockWriter *MockResponseWriter) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWriter := NewMockResponseWriter(ctrl)

			tdkHttpContext := TdkJsonRpcContext{
				Context:              context.Background(),
				Request:              nil,
				Writer:               mockWriter,
				ResponseAlreadyWrite: tt.args.responseAlreadyWrite,
			}
			tt.mock(mockWriter)

			tdkHttpContext.WriteHTTPResponseBatchToJSON(tt.args.response, tt.args.statusCode)
		})
	}
}

func TestJsonRpcContext_ConvertResponseEntityToRpcFormat(t *testing.T) {
	type args struct {
		request  JSONRpcRequestSchema
		response interface{}
	}

	tests := []struct {
		name string
		args args
		want JSONRpcResponseSchema
	}{
		{
			name: "when_given_json_rpc_response_then_expect_write_response",
			args: args{
				response: map[string]interface{}{
					"testing": "testing",
				},
				request: JSONRpcRequestSchema{
					ID: 123,
				},
			},
			want: JSONRpcResponseSchema{
				JSONRpcVersion: "2.0",
				ID:             123,
				Result: map[string]interface{}{
					"testing": "testing",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tdkHttpContext := TdkJsonRpcContext{
				Context: context.Background(),
				Request: nil,
			}

			got := tdkHttpContext.ConvertResponseEntityToRpcFormat(tt.args.request, tt.args.response)
			assert.Equal(t, tt.want, got)
		})
	}
}
