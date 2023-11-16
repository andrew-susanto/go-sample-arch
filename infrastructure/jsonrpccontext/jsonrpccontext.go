package jsonrpccontext

import (
	// golang package
	"context"
	"encoding/json"
	"net/http"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
)

const (
	jsonRpcVersion       = "2.0"
	contentTypeKey       = "Content-Type"
	contentTypeValue     = "application/json"
	internalErrorMessage = "Internal error"
)

//go:generate mockgen -source=jsonrpccontext.go -destination=jsonrpccontext_mock.go -package=jsonrpccontext

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

// TdkJsonRpcContext represents http request context
type TdkJsonRpcContext struct {
	Context              context.Context
	Request              *http.Request
	Writer               ResponseWriter
	ResponseAlreadyWrite bool
}

// JSOnRpcRequestSchema represents json rpc request schema
type JSONRpcRequestSchema struct {
	JSONRpcVersion string          `json:"jsonrpc"`
	MethodName     string          `json:"method"`
	Params         json.RawMessage `json:"params"`
	ID             interface{}     `json:"id"` // id can be null, string, or int
}

// JSONRpcResponseSchema represents json rpc response schema
type JSONRpcResponseSchema struct {
	JSONRpcVersion string              `json:"jsonrpc"`
	Result         interface{}         `json:"result,omitempty"`
	Error          *JSONRpcErrorSchema `json:"error,omitempty"`
	ID             interface{}         `json:"id"` // id can be null, string, or int
}

// JSONRpcErrorSchema represents json rpc error schema
type JSONRpcErrorSchema struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// WriteHTTPResponseToJSON writes response to json using jsonrpc 2.0 spec
func (tdkCtx *TdkJsonRpcContext) WriteHTTPResponseToJSON(response JSONRpcResponseSchema, statusCode int) {
	// prevent multiple write response
	if tdkCtx.ResponseAlreadyWrite {
		return
	}

	js, err := json.Marshal(response)
	if err != nil {
		log.Error(err, nil, "json.Marshal() got error - WriteHTTPResponseToJSON")
		errorResponse := tdkCtx.ConvertResponseErrorToRpcFormat(JSONRpcRequestSchema{}, -32603, internalErrorMessage)

		// best effort to handle error parsing
		js, _ = json.Marshal(errorResponse)
		statusCode = http.StatusInternalServerError
	}

	tdkCtx.Writer.Header().Set("Content-Type", "application/json")
	tdkCtx.Writer.WriteHeader(statusCode)
	tdkCtx.Writer.Write(js)

	tdkCtx.ResponseAlreadyWrite = true
}

// WriteHTTPResponseBatchToJSON writes response to json using jsonrpc 2.0 spec in batch
func (tdkCtx *TdkJsonRpcContext) WriteHTTPResponseBatchToJSON(response []JSONRpcResponseSchema, statusCode int) {
	// prevent multiple write response
	if tdkCtx.ResponseAlreadyWrite {
		return
	}

	js, err := json.Marshal(response)
	if err != nil {
		log.Error(err, nil, "json.Marshal() got error - WriteHTTPResponseBatchToJSON")
		errorResponse := tdkCtx.ConvertResponseErrorToRpcFormat(JSONRpcRequestSchema{}, -32603, internalErrorMessage)

		// best effort to handle error parsing
		js, _ = json.Marshal(errorResponse)
		statusCode = http.StatusInternalServerError
	}

	tdkCtx.Writer.Header().Set(contentTypeKey, contentTypeValue)
	tdkCtx.Writer.WriteHeader(statusCode)
	tdkCtx.Writer.Write(js)

	tdkCtx.ResponseAlreadyWrite = true
}

// ConvertResponseEntityToRpcFormat converts response entity to json rpc response schema
func (tdkCtx *TdkJsonRpcContext) ConvertResponseEntityToRpcFormat(request JSONRpcRequestSchema, response interface{}) JSONRpcResponseSchema {
	resp := JSONRpcResponseSchema{
		JSONRpcVersion: jsonRpcVersion,
		ID:             request.ID,
		Result:         response,
	}

	return resp
}

// ConvertResponseErrorToRpcFormat converts response error to json rpc error response schema
func (tdkCtx *TdkJsonRpcContext) ConvertResponseErrorToRpcFormat(request JSONRpcRequestSchema, errorCode int, message string) JSONRpcResponseSchema {
	resp := JSONRpcResponseSchema{
		JSONRpcVersion: jsonRpcVersion,
		ID:             request.ID,
		Error: &JSONRpcErrorSchema{
			Code:    errorCode,
			Message: message,
		},
	}

	return resp
}
