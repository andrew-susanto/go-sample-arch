package httphandler

import (
	// golang package
	"encoding/json"
	"net/http"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/httpcontext"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
)

// GetUserHandler handles get user request
//
// Returns nil error if success
// Otherwise return non nil error
func (handler *Handler) GetUserHandler(tdkCtx *httpcontext.TdkHttpContext) error {
	var param GetUserParam

	err := json.NewDecoder(tdkCtx.Request.Body).Decode(&param)
	if err != nil {
		log.Error(err, nil, "json.NewDecoder() got error - GetUserHandler")
		return err
	}

	user, err := handler.user.GetUser(tdkCtx.Context, param.ID)
	if err != nil {
		err = errors.Wrap(err).WithCode("HNDL.GUH00")
		log.Error(err, nil, "handler.user.GetUser() got error - GetUserHandler")
		return err
	}

	respUser := GetUserResponse(user)

	tdkCtx.WriteHTTPResponseToJSON(respUser, http.StatusOK)
	return nil
}
