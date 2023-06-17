package config

import (
	// golang package
	"testing"

	// external package
	"github.com/stretchr/testify/assert"
)

func TestConfig_ParseConfig(t *testing.T) {
	type args struct {
		environment string
	}

	tests := []struct {
		name     string
		args     args
		want     Config
		mockFunc func()
	}{
		{
			name: "when_given_environment_dev_expect_config_not_empty",
			args: args{
				environment: "dev",
			},
			mockFunc: func() {
				basePath = "../../"
			},
			want: Config{
				AWS: ConfigAWS{
					Region: "ap-southeast1",
				},
			},
		},
		{
			name: "when_given_environment_staging_expect_config_not_empty",
			args: args{
				environment: "staging",
			},
			mockFunc: func() {
				basePath = "../../"
			},
			want: Config{
				AWS: ConfigAWS{
					Region: "ap-southeast1",
				},
			},
		},
		{
			name: "when_given_environment_production_expect_config_not_empty",
			args: args{
				environment: "production",
			},
			mockFunc: func() {
				basePath = "../../"
			},
			want: Config{
				AWS: ConfigAWS{
					Region: "ap-southeast1",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()
			config := ParseConfig(tt.args.environment)
			assert.Equal(t, tt.want.AWS.Region, config.AWS.Region)
		})
	}
}
