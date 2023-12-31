// Code generated by MockGen. DO NOT EDIT.
// Source: ./repo.go
//
// Generated by this command:
//
//	mockgen -source=./repo.go -destination=./repo_mock.go -package=sqs
//
// Package sqs is a generated GoMock package.
package sqs

import (
	context "context"
	reflect "reflect"

	sqs "github.com/aws/aws-sdk-go-v2/service/sqs"
	gomock "go.uber.org/mock/gomock"
)

// MockSQSPublisher is a mock of SQSPublisher interface.
type MockSQSPublisher struct {
	ctrl     *gomock.Controller
	recorder *MockSQSPublisherMockRecorder
}

// MockSQSPublisherMockRecorder is the mock recorder for MockSQSPublisher.
type MockSQSPublisherMockRecorder struct {
	mock *MockSQSPublisher
}

// NewMockSQSPublisher creates a new mock instance.
func NewMockSQSPublisher(ctrl *gomock.Controller) *MockSQSPublisher {
	mock := &MockSQSPublisher{ctrl: ctrl}
	mock.recorder = &MockSQSPublisherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSQSPublisher) EXPECT() *MockSQSPublisherMockRecorder {
	return m.recorder
}

// SendMessage mocks base method.
func (m *MockSQSPublisher) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendMessage", varargs...)
	ret0, _ := ret[0].(*sqs.SendMessageOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockSQSPublisherMockRecorder) SendMessage(ctx, params any, optFns ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockSQSPublisher)(nil).SendMessage), varargs...)
}

// SendMessageBatch mocks base method.
func (m *MockSQSPublisher) SendMessageBatch(ctx context.Context, params *sqs.SendMessageBatchInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageBatchOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, params}
	for _, a := range optFns {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendMessageBatch", varargs...)
	ret0, _ := ret[0].(*sqs.SendMessageBatchOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMessageBatch indicates an expected call of SendMessageBatch.
func (mr *MockSQSPublisherMockRecorder) SendMessageBatch(ctx, params any, optFns ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, params}, optFns...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessageBatch", reflect.TypeOf((*MockSQSPublisher)(nil).SendMessageBatch), varargs...)
}
