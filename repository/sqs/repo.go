package sqs

import (
	// internal package
	"context"

	// external package
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

//go:generate mockgen -source=./repo.go -destination=./repo_mock.go -package=sqs

type SQSPublisher interface {
	// Delivers a message to the specified queue.
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)

	// You can use SendMessageBatch to send up to 10 messages to the specified queue
	// by assigning either identical or different values to each message (or by not
	// assigning values at all). This is a batch version of SendMessage .
	SendMessageBatch(ctx context.Context, params *sqs.SendMessageBatchInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageBatchOutput, error)
}

type Repository struct {
	sqs SQSPublisher
}

// NewRepository is a function to initialize sqs repository
func NewRepository(sqs SQSPublisher) Repository {
	return Repository{
		sqs: sqs,
	}
}
