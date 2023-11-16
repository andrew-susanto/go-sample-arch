package jsonrpchandler

import (
	// golang package
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"
	"time"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/jsonrpccontext"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/monitor"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/tracer"
)

var (
	contextDeadlineInSeconds = 30
)

type JsonRpcFunc func(tdkCtx *jsonrpccontext.TdkJsonRpcContext, param json.RawMessage) (interface{}, error)

// handleFunc is wrapper for json rpc handler
func (h *Handler) handleFunc(rpcMethodMap map[string]JsonRpcFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancelCtx := context.WithTimeout(r.Context(), time.Duration(contextDeadlineInSeconds)*time.Second)
		defer cancelCtx()

		tdkCtx := jsonrpccontext.TdkJsonRpcContext{
			Context: ctx,
			Request: r,
			Writer:  rw,
		}

		var err error
		var isBatch bool

		body, err := io.ReadAll(r.Body)
		if err != nil {
			err = errors.Wrap(err).WithCode("MDL.HF00")
			log.Error(err, nil, "ioutil.ReadAll() got error - handleFunc")
			response := tdkCtx.ConvertResponseErrorToRpcFormat(jsonrpccontext.JSONRpcRequestSchema{}, -32700, "Parse error")
			tdkCtx.WriteHTTPResponseToJSON(response, 500)
			return
		}

		// determine if schema is batch or single
		var rpcSchema interface{}
		err = json.Unmarshal(body, &rpcSchema)
		if err != nil {
			err = errors.Wrap(err).WithCode("MDL.HF01")
			log.Error(err, nil, "json.NewDecoder() got error - handleFunc")
			response := tdkCtx.ConvertResponseErrorToRpcFormat(jsonrpccontext.JSONRpcRequestSchema{}, -32700, "Parse error")
			tdkCtx.WriteHTTPResponseToJSON(response, 500)
			return
		}

		switch rpcSchema.(type) {
		case []interface{}:
			isBatch = true
		case map[string]interface{}:
			isBatch = false
		default:
			err = errors.New("Invalid rpcschema").WithCode("MDL.HF02")
			log.Error(err, nil, "json.NewDecoder() got error - handleFunc")
			response := tdkCtx.ConvertResponseErrorToRpcFormat(jsonrpccontext.JSONRpcRequestSchema{}, -32700, "Parse error")
			tdkCtx.WriteHTTPResponseToJSON(response, 500)
			return
		}

		if isBatch {
			var requestBatch []json.RawMessage
			err = json.Unmarshal(body, &requestBatch)
			if err != nil {
				err = errors.Wrap(err).WithCode("MDL.HF03")
				log.Error(err, nil, "json.NewDecoder() got error - handleFunc")
				response := tdkCtx.ConvertResponseErrorToRpcFormat(jsonrpccontext.JSONRpcRequestSchema{}, -32700, "Parse error")
				tdkCtx.WriteHTTPResponseToJSON(response, 500)
				return
			}

			if len(requestBatch) == 0 {
				err = errors.New("Invalid request empty array").WithCode("MDL.HF04")
				log.Error(err, nil, "invalid request empty array - handleFunc")
				response := tdkCtx.ConvertResponseErrorToRpcFormat(jsonrpccontext.JSONRpcRequestSchema{}, -32600, "Invalid Request")
				tdkCtx.WriteHTTPResponseToJSON(response, 500)
				return
			}

			parsedRequestBatch := make([]jsonrpccontext.JSONRpcRequestSchema, len(requestBatch))
			isValidRequestBatch := make([]bool, len(requestBatch))
			var responseBatchCount int

			for idx, request := range requestBatch {
				requestSingle, isValid := h.validateRpcRequest(request)

				parsedRequestBatch[idx] = requestSingle
				isValidRequestBatch[idx] = isValid

				// count required response
				if !isValid || requestSingle.ID != nil {
					responseBatchCount = responseBatchCount + 1
				}
			}

			responseBatch := make([]jsonrpccontext.JSONRpcResponseSchema, responseBatchCount)

			var wg sync.WaitGroup
			var responseBatchIdx int
			for idx, request := range parsedRequestBatch {
				// invalid request
				if !isValidRequestBatch[idx] {
					responseBatch[responseBatchIdx] = tdkCtx.ConvertResponseErrorToRpcFormat(request, -32600, "Invalid Request")
					responseBatchIdx = responseBatchIdx + 1
					continue
				}

				// notification
				if request.ID == nil {
					wg.Add(1)
					go func(request jsonrpccontext.JSONRpcRequestSchema) {
						h.handleRpcFunc(&tdkCtx, rpcMethodMap, request)
						wg.Done()
					}(request)
					continue
				}

				// non-notification therefore expect response
				wg.Add(1)
				go func(responseBatchIdx int, request jsonrpccontext.JSONRpcRequestSchema) {
					responseBatch[responseBatchIdx] = h.handleRpcFunc(&tdkCtx, rpcMethodMap, request)
					wg.Done()
				}(responseBatchIdx, request)
				responseBatchIdx = responseBatchIdx + 1
			}

			wg.Wait()
			if len(responseBatch) > 0 {
				tdkCtx.WriteHTTPResponseBatchToJSON(responseBatch, 200)
			}
			return
		}

		// non-batch
		request, isValid := h.validateRpcRequest(body)
		if !isValid {
			response := tdkCtx.ConvertResponseErrorToRpcFormat(request, -32600, "Invalid Request")
			tdkCtx.WriteHTTPResponseToJSON(response, 500)
			return
		}

		response := h.handleRpcFunc(&tdkCtx, rpcMethodMap, request)
		if request.ID != nil {
			tdkCtx.WriteHTTPResponseToJSON(response, 200)
		}
	}
}

// validateRpcRequest validates rpc request
//
// Returns parsed request and true if success
// Otherwise return empty request and false
func (h *Handler) validateRpcRequest(request json.RawMessage) (jsonrpccontext.JSONRpcRequestSchema, bool) {
	var requestParsed jsonrpccontext.JSONRpcRequestSchema
	var err error

	err = json.Unmarshal(request, &requestParsed)
	if err != nil {
		err = errors.Wrap(err).WithCode("MDL.VRR00")
		log.Error(err, nil, "json.NewDecoder() got error - validateRpcRequest")
		return requestParsed, false
	}

	if requestParsed.JSONRpcVersion != "2.0" {
		err = errors.New("Invalid RPC version").WithCode("MDL.VRR01")
		log.Error(err, nil, "invalid rpc version - validateRpcRequest")
		return requestParsed, false
	}

	if requestParsed.MethodName == "" {
		err = errors.New("Invalid method name nil").WithCode("MDL.VRR02")
		log.Error(err, nil, "invalid method name nil - validateRpcRequest")
		return requestParsed, false
	}

	if reflect.TypeOf(requestParsed.MethodName).Kind() != reflect.String {
		err = errors.New("Invalid method name type").WithCode("MDL.VRR03")
		log.Error(err, nil, "invalid method name type - validateRpcRequest")
		return requestParsed, false
	}

	return requestParsed, true
}

// handleRpcFunc is inner handler for rpc to handle single request
func (h *Handler) handleRpcFunc(tdkCtx *jsonrpccontext.TdkJsonRpcContext, rpcMethodMap map[string]JsonRpcFunc, request jsonrpccontext.JSONRpcRequestSchema) jsonrpccontext.JSONRpcResponseSchema {
	metricsTag := []string{}

	var isSuccess bool
	var response jsonrpccontext.JSONRpcResponseSchema

	start := time.Now()
	fnResponse, err := h.handleExecuteRpcFunc(tdkCtx, rpcMethodMap, request)
	if err != nil {
		var errorConverted errors.Error
		switch errConvert := err.(type) {
		case errors.Error:
			isSuccess = errConvert.EType == errors.USER
			errorConverted = errConvert
		case error:
			errorConverted = errors.Wrap(errConvert).WithCode("MDL.HRF00").WithNumber(-32000)
		}

		response = tdkCtx.ConvertResponseErrorToRpcFormat(request, errorConverted.ENumber, err.Error())
		metricsTag = append(metricsTag, fmt.Sprintf("errorcode:%v", errorConverted.ECode))
		metricsTag = append(metricsTag, fmt.Sprintf("errornumber:%v", errorConverted.ENumber))
	} else {
		isSuccess = true
		response = tdkCtx.ConvertResponseEntityToRpcFormat(request, fnResponse)
	}

	// push metrics
	countMetricsName := fmt.Sprintf("%s.count", monitor.MetricsPrefix)
	metricsTag = append(metricsTag, fmt.Sprintf("success:%v", isSuccess))
	metricsTag = append(metricsTag, fmt.Sprintf("method:%v", request.MethodName))
	h.monitor.Incr(countMetricsName, metricsTag, 1)

	gaugeMetricsName := fmt.Sprintf("%s.duration", monitor.MetricsPrefix)
	h.monitor.Gauge(gaugeMetricsName, float64(time.Since(start).Microseconds()), []string{
		fmt.Sprintf("method:%v", request.MethodName),
	}, 1)

	return response
}

// handleExecuteRpcFunc is inner handler for rpc to execute single request
func (h *Handler) handleExecuteRpcFunc(tdkCtx *jsonrpccontext.TdkJsonRpcContext, rpcMethodMap map[string]JsonRpcFunc, request jsonrpccontext.JSONRpcRequestSchema) (interface{}, error) {
	var err error

	// recover from panic to prevent server crash
	defer func() {
		r := recover()
		if r != nil {
			switch t := r.(type) {
			case error:
				err = errors.Wrap(err).WithCode("MDL.HERF00").WithNumber(-32000)
			default:
				err = errors.New(fmt.Sprintf("%v", t)).WithCode("MDL.HERF01").WithNumber(-32000)
			}
		}
	}()

	fn, methodExists := rpcMethodMap[request.MethodName]
	if !methodExists {
		err = errors.New("Method not found").WithCode("MDL.HERF02").WithNumber(-32601)
		log.Error(err, nil, "method not found rpc - handleExecuteRpcFunc")
		return nil, err
	}

	spanName := fmt.Sprintf("%s - %s - %s", tdkCtx.Request.URL.Path, tdkCtx.Request.Method, request.MethodName)
	ctx, span := tracer.Start(tdkCtx.Context, spanName)
	defer span.End()

	// replace context with otel tracer context
	tdkCtx.Context = ctx

	// run func if no error
	responseFn, err := fn(tdkCtx, request.Params)
	if err != nil {
		err = errors.Wrap(err).WithCode("MDL.HERF03").WithNumber(-32000)
		log.Error(err, nil, "fn() got error - handleExecuteRpcFunc")
		return nil, err
	}

	return responseFn, nil
}
