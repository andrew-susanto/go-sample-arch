package infrastructure

import (
	// golang package
	"testing"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/config"

	// external package
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestInfrastructure_InitInfrastructure(t *testing.T) {
	type args struct {
		config config.Config
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(*MockParamStore) Infrastructure
	}{
		{
			name: "when_given_config_then_return_infrastructure",
			args: args{
				config: config.Config{
					FeatureFlag: map[string]interface{}{
						"test": "test",
					},
				},
			},
			mockFunc: func(mock *MockParamStore) Infrastructure {
				return Infrastructure{
					Config: config.Config{
						FeatureFlag: map[string]interface{}{
							"test": "test",
						},
					},
					paramstore: mock,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockParamStore := NewMockParamStore(ctrl)
			want := tt.mockFunc(mockParamStore)
			got := InitInfrastructure(tt.args.config, mockParamStore)

			assert.Equal(t, want, got)
		})
	}
}

func TestInfrastructure_GetParamStoreValue(t *testing.T) {
	type args struct {
		key string
	}

	tests := []struct {
		name     string
		args     args
		want     string
		mockFunc func(*MockParamStore)
	}{
		{
			name: "when_given_param_store_key_then_return_value",
			args: args{
				key: "testing",
			},
			want: "value_testing",
			mockFunc: func(mock *MockParamStore) {
				mock.EXPECT().GetValue("testing").Return("value_testing")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockParamStore := NewMockParamStore(ctrl)
			infrastructure := InitInfrastructure(config.Config{}, mockParamStore)

			tt.mockFunc(mockParamStore)

			got := infrastructure.GetParamStoreValue(tt.args.key)
			assert.Equal(t, tt.want, got)
		})
	}
}
