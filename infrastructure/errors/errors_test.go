package errors

import (
	// golang package
	"errors"
	"testing"

	// external package
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestError_Wrap(t *testing.T) {
	type args struct {
		e error
	}

	tests := []struct {
		name string
		args args
		want Error
	}{
		{
			name: "when_given_go_error_expect_wrap_error",
			args: args{
				e: errors.New("test"),
			},
			want: Error{
				E:     errors.New("test"),
				EType: SYSTEM,
			},
		},
		{
			name: "when_given_nil_expect_unknown_error",
			args: args{
				e: nil,
			},
			want: Error{
				E:     errors.New("unknown error"),
				EType: SYSTEM,
				ECode: "INF.ERR00",
			},
		},
		{
			name: "when_given_custom_error_expect_do_nothing",
			args: args{
				e: Error{
					E:     errors.New("custom errorr"),
					EType: USER,
					ECode: "INF.ERR01",
				},
			},
			want: Error{
				E:     errors.New("custom errorr"),
				EType: USER,
				ECode: "INF.ERR01",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			got := Wrap(tt.args.e)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestError_WithType(t *testing.T) {
	type args struct {
		e     Error
		eType int
	}

	tests := []struct {
		name string
		args args
		want Error
	}{
		{
			name: "when_given_error_type_system_expect_error_type_updated",
			args: args{
				e: Error{
					E: errors.New("test"),
				},
				eType: SYSTEM,
			},
			want: Error{
				E:     errors.New("test"),
				EType: SYSTEM,
			},
		},
		{
			name: "when_given_error_type_user_expect_error_type_updated",
			args: args{
				e: Error{
					E: errors.New("test"),
				},
				eType: USER,
			},
			want: Error{
				E:     errors.New("test"),
				EType: USER,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			got := tt.args.e.WithType(tt.args.eType)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestError_WithCode(t *testing.T) {
	type args struct {
		e     Error
		eCode string
	}

	tests := []struct {
		name string
		args args
		want Error
	}{
		{
			name: "when_given_error_code_and_current_error_code_empty_expect_error_code_updated",
			args: args{
				e: Error{
					E: errors.New("test"),
				},
				eCode: "INF.ERR01",
			},
			want: Error{
				E:     errors.New("test"),
				ECode: "INF.ERR01",
			},
		},
		{
			name: "when_given_error_code_and_current_error_code_not_empty_expect_error_code_not_updated",
			args: args{
				e: Error{
					E:     errors.New("test"),
					ECode: "INF.ERR00",
				},
				eCode: "INF.ERR01",
			},
			want: Error{
				E:     errors.New("test"),
				ECode: "INF.ERR00",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			got := tt.args.e.WithCode(tt.args.eCode)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestError_WithNumber(t *testing.T) {
	type args struct {
		e       Error
		eNumber int
	}

	tests := []struct {
		name string
		args args
		want Error
	}{
		{
			name: "when_given_error_number_expect_error_number_updated",
			args: args{
				e: Error{
					E:       errors.New("test"),
					ENumber: -32000,
				},
				eNumber: -31000,
			},
			want: Error{
				E:       errors.New("test"),
				ENumber: -31000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			got := tt.args.e.WithNumber(tt.args.eNumber)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestError_Error(t *testing.T) {
	type args struct {
		e Error
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "when_given_error_expect_return_error_string",
			args: args{
				e: Error{
					E:       errors.New("testing_error_string"),
					ENumber: -32000,
				},
			},
			want: "testing_error_string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			got := tt.args.e.Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestError_Is(t *testing.T) {
	err1 := Error{
		E: errors.New("test"),
	}
	err2 := Error{
		E: errors.New("test"),
	}

	// assert same errror
	got := Is(err1, err1)
	assert.Equal(t, true, got)

	// assert different error
	got = Is(err1, err2)
	assert.Equal(t, false, got)
}
