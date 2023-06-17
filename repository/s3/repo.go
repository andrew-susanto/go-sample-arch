package s3

import (
	// golang package
	"context"

	// external package
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//go:generate mockgen -source=./repo.go -destination=./repo_mock.go -package=s3

type S3Service interface {
	// Retrieves objects from Amazon S3.
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)

	// Adds an object to a bucket.
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)

	// Returns some or all (up to 1,000) of the objects in a bucket with each request.
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

type Repository struct {
	s3 S3Service
}

// NewRepository is a function to initialize s3 repository
func NewRepository(s3service S3Service) Repository {
	return Repository{
		s3: s3service,
	}
}
