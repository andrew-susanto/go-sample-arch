package jsonrpchandler

import (
	// golang package
	"encoding/json"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/jsonrpccontext"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/tracer"
)

// GetUserHandler gets user by given id
//
// Return user and nil error when success
// Otherwise return empty user and non nil error
func (handler *Handler) GetUserHandler(tdkCtx *jsonrpccontext.TdkJsonRpcContext, params json.RawMessage) (interface{}, error) {
	ctx, span := tracer.Start(tdkCtx.Context, "handler.jsonrpchandler.GetUserHandler")
	defer span.End()

	var param GetUserParam

	err := json.Unmarshal(params, &param)
	if err != nil {
		err = errors.Wrap(err).WithCode("HNDL.GUH00")
		log.Error(err, nil, "json.Umarhsal() got error - GetUserHandler")
		return nil, err
	}

	user, err := handler.user.GetUser(ctx, param.ID)
	if err != nil {
		err = errors.Wrap(err).WithCode("HNDL.GUH00")
		log.Error(err, nil, "handler.user.GetUser() got error - GetUserHandler")
		return nil, err
	}

	respUser := GetUserResponse(user)
	return respUser, nil
}
