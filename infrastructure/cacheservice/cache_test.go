package cacheservice

import (
	// golang package
	"testing"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/secretmanager"

	// external package
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestInfrastructure_InitCache(t *testing.T) {
	type args struct {
		config secretmanager.SecretsRedis
	}

	tests := []struct {
		name    string
		args    args
		want    Cache
		wantErr error
	}{
		{
			name: "when_given_host_port_and_password_then_return_cache_service",
			args: args{
				config: secretmanager.SecretsRedis{
					Host:     "localhost",
					Port:     "1234",
					Password: "test",
				},
			},
			want: redis.NewClient(&redis.Options{
				Addr:     "localhost:1234",
				Password: "test",
				DB:       0,
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := InitCache(tt.args.config)
			_, ok := client.(*redis.Client)
			assert.True(t, ok)
		})
	}
}
