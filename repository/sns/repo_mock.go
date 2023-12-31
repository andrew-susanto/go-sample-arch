// Code generated by MockGen. DO NOT EDIT.
// Source: ./repo.go
//
// Generated by this command:
//
//	mockgen -source=./repo.go -destination=./repo_mock.go -package=sns
//
// Package sns is a generated GoMock package.
package sns

import (
	context "context"
	reflect "reflect"

	sns "github.com/aws/aws-sdk-go-v2/service/sns"
	gomock "go.uber.org/mock/gomock"
)

// MockSNSPublisher is a mock of SNSPublisher interface.
type MockSNSPublisher struct {
	ctrl     *gomock.Controller
	recorder *MockSNSPublisherMockRecorder
}

// MockSNSPublisherMockRecorder is the mock recorder for MockSNSPublisher.
type MockSNSPublisherMockRecorder struct {
	mock *MockSNSPublisher
}

// NewMockSNSPublisher creates a new mock instance.
func NewMockSNSPublisher(ctrl *gomock.Controller) *MockSNSPublisher {
	mock := &MockSNSPublisher{ctrl: ctrl}
	mock.recorder = &MockSNSPublisherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSNSPublisher) EXPECT() *MockSNSPublisherMockRecorder {
	return m.recorder
}

// Publish mocks base method.
func (m *MockSNSPublisher) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Publish", varargs...)
	ret0, _ := ret[0].(*sns.PublishOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Publish indicates an expected call of Publish.
func (mr *MockSNSPublisherMockRecorder) Publish(ctx, params any, optFns ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockSNSPublisher)(nil).Publish), varargs...)
}
