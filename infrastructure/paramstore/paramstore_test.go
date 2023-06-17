package paramstore

import (
	// golang package
	"errors"
	"testing"

	// external package
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
)

func TestParamStore_InitParamStore(t *testing.T) {
	config := aws.Config{}

	got := InitParamstore(config)
	assert.NotNil(t, got)
}

func TestParamStore_GetValue(t *testing.T) {
	type args struct {
		key string
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(*MockParamStoreClientInterface)
		want     string
	}{
		{
			name: "when_given_key_then_return_parameter_value",
			args: args{
				key: "key_testing",
			},
			mockFunc: func(mock *MockParamStoreClientInterface) {
				mock.EXPECT().GetParameter(gomock.Any(), &ssm.GetParameterInput{
					Name: aws.String("key_testing"),
				}).Return(&ssm.GetParameterOutput{
					Parameter: &types.Parameter{
						Value: aws.String("value_testing"),
					},
				}, nil)
			},
			want: "value_testing",
		},
		{
			name: "when_given_key_and_parameter_value_nil_then_return_empty_string",
			args: args{
				key: "key_testing",
			},
			mockFunc: func(mock *MockParamStoreClientInterface) {
				mock.EXPECT().GetParameter(gomock.Any(), &ssm.GetParameterInput{
					Name: aws.String("key_testing"),
				}).Return(&ssm.GetParameterOutput{
					Parameter: &types.Parameter{
						Value: nil,
					},
				}, nil)
			},
			want: "",
		},
		{
			name: "when_given_key_and_get_parameter_failed_then_return_empty_string",
			args: args{
				key: "key_testing",
			},
			mockFunc: func(mock *MockParamStoreClientInterface) {
				mock.EXPECT().GetParameter(gomock.Any(), &ssm.GetParameterInput{
					Name: aws.String("key_testing"),
				}).Return(&ssm.GetParameterOutput{
					Parameter: &types.Parameter{
						Value: aws.String("value_testing"),
					},
				}, errors.New("test_error"))
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockParamStore := NewMockParamStoreClientInterface(ctrl)
			tt.mockFunc(mockParamStore)

			paramStore := ParamStore{
				client: mockParamStore,
			}

			got := paramStore.GetValue(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}
