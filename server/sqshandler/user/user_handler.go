package user

import (
	// golang package
	"context"
	"encoding/json"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/tracer"

	// external package
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// GetUserParam is the parameter for GetUserHandler
//
// Returns nil error if success
// Otherwise return non nil error
func (handler *Handler) GetUserHandler(ctx context.Context, message types.Message) error {
	ctx, span := tracer.Start(ctx, "handler.sqshandler.GetUserHandler")
	defer span.End()

	var param GetUserParam

	err := json.Unmarshal([]byte(*message.Body), &param)
	if err != nil {
		err = errors.Wrap(err).WithCode("HNDL.GUH00")
		log.Error(err, nil, "json.Umarhsal() got error - GetUserHandler")
		return err
	}

	_, err = handler.user.GetUser(ctx, param.ID)
	if err != nil {
		err = errors.Wrap(err).WithCode("HNDL.GUH00")
		log.Error(err, nil, "handler.user.GetUser() got error - GetUserHandler")
		return err
	}

	return nil
}
