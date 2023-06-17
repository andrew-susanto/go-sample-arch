package sns

import (
	// internal package
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

//go:generate mockgen -source=./repo.go -destination=./repo_mock.go -package=sns

type SNSPublisher interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

type Repository struct {
	sns SNSPublisher
}

// NewRepository is a function to initialize sns repository
func NewRepository(sns SNSPublisher) Repository {
	return Repository{
		sns: sns,
	}
}
