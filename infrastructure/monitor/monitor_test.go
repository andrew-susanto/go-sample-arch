package monitor

import (
	// golang package
	"testing"

	// external package
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMonitor_InitMonitor(t *testing.T) {
	type args struct {
		environment string
		serviceName string
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "when_given_param_then_return_monitor",
			args: args{
				environment: "dev",
				serviceName: "cxpcrmapp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			got := InitMonitor(tt.args.environment, tt.args.serviceName)
			assert.NotNil(t, got)
		})
	}
}
