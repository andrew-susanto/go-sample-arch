// Code generated by MockGen. DO NOT EDIT.
// Source: ./repo.go
//
// Generated by this command:
//
//	mockgen -source=./repo.go -destination=./repo_mock.go -package=cache
//
// Package cache is a generated GoMock package.
package cache

import (
	context "context"
	reflect "reflect"
	time "time"

	redis "github.com/redis/go-redis/v9"
	gomock "go.uber.org/mock/gomock"
)

// MockRedis is a mock of Redis interface.
type MockRedis struct {
	ctrl     *gomock.Controller
	recorder *MockRedisMockRecorder
}

// MockRedisMockRecorder is the mock recorder for MockRedis.
type MockRedisMockRecorder struct {
	mock *MockRedis
}

// NewMockRedis creates a new mock instance.
func NewMockRedis(ctrl *gomock.Controller) *MockRedis {
	mock := &MockRedis{ctrl: ctrl}
	mock.recorder = &MockRedisMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedis) EXPECT() *MockRedisMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockRedis) Get(ctx context.Context, key string) *redis.StringCmd {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(*redis.StringCmd)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockRedisMockRecorder) Get(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRedis)(nil).Get), ctx, key)
}

// Set mocks base method.
func (m *MockRedis) Set(ctx context.Context, key string, value any, duration time.Duration) *redis.StatusCmd {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, value, duration)
	ret0, _ := ret[0].(*redis.StatusCmd)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockRedisMockRecorder) Set(ctx, key, value, duration any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockRedis)(nil).Set), ctx, key, value, duration)
}