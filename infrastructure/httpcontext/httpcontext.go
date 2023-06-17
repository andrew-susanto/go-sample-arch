package httpcontext

import (
	// golang package
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
)

const (
	contentTypeKey   = "Content-Type"
	contentTypeValue = "application/json"

	badRequestErrorMessage = "BAD_REQUEST"
	serverErrorMessage     = "SERVER_ERROR"

	badRequestStatusCode  = 400
	serverErrorStatusCode = 500
)

//go:generate mockgen -source=httpcontext.go -destination=httpcontext_mock.go -package=httpcontext

type ResponseWriter interface {
	// Header returns the header map that will be sent by
	// WriteHeader. The Header map also is the mechanism with which
	// Handlers can set HTTP trailers.
	Header() http.Header

	// Write writes the data to the connection as part of an HTTP reply.
	Write([]byte) (int, error)

	// WriteHeader sends an HTTP response header with the provided
	// status code.
	WriteHeader(statusCode int)
}

// TdkHttpContext represents http request context
type TdkHttpContext struct {
	Context              context.Context
	Request              *http.Request
	Writer               ResponseWriter
	ResponseAlreadyWrite bool
	ResponseUserContext  interface{}
}

// WriteHTTPResponseToJSON write http response with given data and status code
// data must be marshall-able value and status code must be a valid status code
// otherwise internal server error will be returned
func (tdkCtx *TdkHttpContext) WriteHTTPResponseToJSON(responseJson interface{}, statusCode int) {
	// prevent multiple write response
	if tdkCtx.ResponseAlreadyWrite {
		return
	}

	resp := map[string]interface{}{
		"data":        responseJson,
		"userContext": tdkCtx.ResponseUserContext,
	}

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(tdkCtx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	tdkCtx.Writer.Header().Set(contentTypeKey, contentTypeValue)
	tdkCtx.Writer.WriteHeader(statusCode)
	tdkCtx.Writer.Write(js)

	tdkCtx.ResponseAlreadyWrite = true
}

// WriteHTTPResponseErrorToJSON writes http response with given error
func (tdkCtx *TdkHttpContext) WriteHTTPResponseErrorToJSON(err error) {
	// prevent multiple write response
	if tdkCtx.ResponseAlreadyWrite {
		return
	}

	statusCode := serverErrorStatusCode
	errorType := serverErrorMessage

	switch errConvert := err.(type) {
	case errors.Error:
		if errConvert.EType == errors.USER {
			statusCode = badRequestStatusCode
			errorType = badRequestErrorMessage
		}
	}

	resp := map[string]interface{}{
		"errorType":        errorType,
		"userErrorMessage": err.Error(),
		"errorMessage":     err.Error(),
		"statusCode":       strconv.FormatInt(int64(statusCode), 10),
	}

	js, err := json.Marshal(resp)
	if err != nil {
		http.Error(tdkCtx.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	tdkCtx.Writer.Header().Set(contentTypeKey, contentTypeValue)
	tdkCtx.Writer.WriteHeader(statusCode)
	tdkCtx.Writer.Write(js)

	tdkCtx.ResponseAlreadyWrite = true
}

// SetResponseUserContext set user context on the response
// this function must be called before write http response
// otherwise this function will have no effect
func (tdkCtx *TdkHttpContext) SetResponseUserContext(userContext interface{}) {
	tdkCtx.ResponseUserContext = userContext
}
