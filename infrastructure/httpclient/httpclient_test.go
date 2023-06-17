package httpclient

import (
	// golang package
	"net/http"
	"testing"
	"time"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/config"

	// external package
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/mock/gomock"
)

func TestHttpClient_InitHttpClient(t *testing.T) {
	type args struct {
		config config.Config
	}

	tests := []struct {
		name string
		args args
		want HttpClient
	}{
		{
			name: "when_given_config_then_return_httpclient",
			args: args{
				config: config.Config{
					HttpClient: map[string]config.ConfigHttpClient{
						"test": {
							Timeout:                  30,
							MaxConcurrentRequest:     10,
							ErrorPercentageThreshold: 25,
						},
					},
				},
			},
			want: HttpClient{
				http: &http.Client{
					Timeout:   time.Duration(maxHttpClientTimeoutInSeconds) * time.Second,
					Transport: otelhttp.NewTransport(http.DefaultTransport),
				},
				config: config.Config{
					HttpClient: map[string]config.ConfigHttpClient{
						"test": {
							Timeout:                  30,
							MaxConcurrentRequest:     10,
							ErrorPercentageThreshold: 25,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			got := InitHttpClient(tt.args.config)
			assert.Equal(t, tt.want.config, got.config)
		})
	}
}

func TestHttpClient_Wrap(t *testing.T) {
	type args struct {
		clientName string
		config     config.Config
	}

	tests := []struct {
		name string
		args args
		want error
	}{
		{
			name: "when_given_client_name_exists_expect_use_configured_client_name",
			args: args{
				config: config.Config{
					HttpClient: map[string]config.ConfigHttpClient{
						"test": {
							Timeout:                  30,
							MaxConcurrentRequest:     10,
							ErrorPercentageThreshold: 100,
						},
					},
				},
				clientName: "test",
			},
			want: nil,
		},
		{
			name: "when_given_client_name_not_exists_expect_use_configured_client_name",
			args: args{
				config: config.Config{
					HttpClient: map[string]config.ConfigHttpClient{
						"test": {
							Timeout:                  30,
							MaxConcurrentRequest:     10,
							ErrorPercentageThreshold: 100,
						},
					},
				},
				clientName: "test_different",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			httpClient := InitHttpClient(tt.args.config)
			got := httpClient.Wrap(tt.args.clientName, func() error { return nil })

			assert.Equal(t, tt.want, got)
		})
	}
}
