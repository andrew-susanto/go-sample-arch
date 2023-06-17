package httpcontext

import (
	// golang package
	"context"
	"net/http"
	"testing"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"

	// external package
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHttpContext_WriteHTTPResponseToJSON(t *testing.T) {
	type args struct {
		responseJson         interface{}
		statusCode           int
		responseAlreadyWrite bool
		responseUserContext  interface{}
	}

	tests := []struct {
		name string
		args args
		mock func(mockWriter *MockResponseWriter)
	}{
		{
			name: "when_given_response_and_status_code_then_expect_write_response_to_http_writer",
			args: args{
				responseJson: map[string]interface{}{
					"testing": "testing_response",
				},
				statusCode: 200,
			},
			mock: func(mockWriter *MockResponseWriter) {
				mockWriter.EXPECT().Header().Return(http.Header{})
				mockWriter.EXPECT().WriteHeader(200)
				mockWriter.EXPECT().Write([]byte("{\"data\":{\"testing\":\"testing_response\"},\"userContext\":null}"))
			},
		},
		{
			name: "when_given_invalid_response_then_expect_write_error",
			args: args{
				responseJson: func() {},
				statusCode:   200,
			},
			mock: func(mockWriter *MockResponseWriter) {
				mockWriter.EXPECT().Header().Return(http.Header{}).Times(2)
				mockWriter.EXPECT().WriteHeader(500)
				mockWriter.EXPECT().Write([]byte("json: unsupported type: func()\n"))
			},
		},
		{
			name: "when_response_already_write_then_expect_no_write",
			args: args{
				responseJson:         func() {},
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

			tdkHttpContext := TdkHttpContext{
				Context:              context.Background(),
				Request:              nil,
				Writer:               mockWriter,
				ResponseAlreadyWrite: tt.args.responseAlreadyWrite,
				ResponseUserContext:  tt.args.responseUserContext,
			}
			tt.mock(mockWriter)

			tdkHttpContext.WriteHTTPResponseToJSON(tt.args.responseJson, tt.args.statusCode)
		})
	}
}

func TestHttpContext_WriteHTTPResponseErrorToJSON(t *testing.T) {
	type args struct {
		error                error
		responseAlreadyWrite bool
		responseUserContext  interface{}
	}

	tests := []struct {
		name string
		args args
		mock func(mockWriter *MockResponseWriter)
	}{
		{
			name: "when_response_already_write_then_expect_no_write",
			args: args{
				error:                errors.New("testing"),
				responseAlreadyWrite: true,
			},
			mock: func(mockWriter *MockResponseWriter) {},
		},
		{
			name: "wheen_given_server_error_then_expect_write_server_error",
			args: args{
				error:                errors.New("testing").WithType(errors.SYSTEM).WithCode("UC.000"),
				responseAlreadyWrite: false,
			},
			mock: func(mockWriter *MockResponseWriter) {
				mockWriter.EXPECT().Header().Return(http.Header{})
				mockWriter.EXPECT().WriteHeader(500)
				mockWriter.EXPECT().Write([]byte("{\"errorMessage\":\"testing\",\"errorType\":\"SERVER_ERROR\",\"statusCode\":\"500\",\"userErrorMessage\":\"testing\"}"))
			},
		},
		{
			name: "wheen_given_bad_request_error_then_expect_write_bad_request_error",
			args: args{
				error:                errors.New("testing").WithType(errors.USER).WithCode("UC.000"),
				responseAlreadyWrite: false,
			},
			mock: func(mockWriter *MockResponseWriter) {
				mockWriter.EXPECT().Header().Return(http.Header{})
				mockWriter.EXPECT().WriteHeader(400)
				mockWriter.EXPECT().Write([]byte("{\"errorMessage\":\"testing\",\"errorType\":\"BAD_REQUEST\",\"statusCode\":\"400\",\"userErrorMessage\":\"testing\"}"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockWriter := NewMockResponseWriter(ctrl)

			tdkHttpContext := TdkHttpContext{
				Context:              context.Background(),
				Request:              nil,
				Writer:               mockWriter,
				ResponseAlreadyWrite: tt.args.responseAlreadyWrite,
				ResponseUserContext:  tt.args.responseUserContext,
			}
			tt.mock(mockWriter)

			tdkHttpContext.WriteHTTPResponseErrorToJSON(tt.args.error)
		})
	}
}

func TestHttpContext_SetResponseUserContext(t *testing.T) {
	type args struct {
		userContext interface{}
	}

	tests := []struct {
		name            string
		args            args
		wantUserContext interface{}
	}{
		{
			name: "when_given_user_context_then_expect_user_context_updated",
			args: args{
				userContext: map[string]interface{}{
					"testing": "testing",
				},
			},
			wantUserContext: map[string]interface{}{
				"testing": "testing",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			tdkHttpContext := TdkHttpContext{
				Context: context.Background(),
				Request: nil,
			}

			tdkHttpContext.SetResponseUserContext(tt.args.userContext)
			assert.Equal(t, tt.wantUserContext, tdkHttpContext.ResponseUserContext)
		})
	}
}
