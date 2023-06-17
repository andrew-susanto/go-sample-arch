package cacheservice

import (
	// golang package
	"testing"

	// external package
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestInfrastructure_InitCache(t *testing.T) {
	type args struct {
		host     string
		port     string
		password string
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
				host:     "localhost",
				port:     "1234",
				password: "test",
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
			client := InitCache(tt.args.host, tt.args.port, tt.args.password)
			_, ok := client.(*redis.Client)
			assert.True(t, ok)
		})
	}
}
