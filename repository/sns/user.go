package sns

import (
	// golang package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"

	// external package
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// GetUserByID gets user by given id from sns
//
// Return user and nil error when success
// Otherwise return empty user and non nil error
func (repo *Repository) GetUserByID(ctx context.Context, id int64) (User, error) {
	_, err := repo.sns.Publish(ctx, &sns.PublishInput{})
	if err != nil {
		err = errors.Wrap(err).WithCode("RPST.GUBI00")
		log.Error(err, id, "repo.sqs.SendMessage() got error - GetUserByID")
		return User{}, err
	}

	return User{
		ID:     id,
		Name:   "dummy user name",
		Gender: 3,
	}, nil
}
