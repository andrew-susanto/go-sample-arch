package httphandler

import (
	// golang package
	"context"
	"fmt"
	"net/http"
	"time"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/httpcontext"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/monitor"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/tracer"
)

var (
	contextDeadlineInSeconds = 30
)

// handleFunc is a wrapper for http handler function
func (h *Handler) handleFunc(fn func(tdkCtx *httpcontext.TdkHttpContext) error) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// recover from panic to prevent server crash
		defer func() {
			r := recover()
			if r != nil {
				var err error
				switch t := r.(type) {
				case string:
					err = errors.New(t).WithCode("MDL.HF00")
				case error:
					err = errors.Wrap(err).WithCode("MDL.HF01")
				default:
					err = errors.New("Unknown error").WithCode("MDL.HF02")
				}
			}
		}()

		ctx, cancelCtx := context.WithTimeout(r.Context(), time.Duration(contextDeadlineInSeconds)*time.Second)
		defer cancelCtx()

		// replace context with otel tracer context
		spanName := fmt.Sprintf("%s - %s", r.URL.Path, r.Method)
		ctx, span := tracer.Start(ctx, spanName)
		defer span.End()

		start := time.Now()
		httpCtx := httpcontext.TdkHttpContext{
			Context: ctx,
			Request: r,
			Writer:  rw,
		}
		err := fn(&httpCtx)
		duration := time.Since(start)

		// determine if request is success
		var isSuccess bool
		var errorConverted errors.Error

		switch errConvert := err.(type) {
		case errors.Error:
			isSuccess = errConvert.EType == errors.USER
			errorConverted = errConvert
		case error:
			isSuccess = errConvert == nil
			errorConverted = errors.Wrap(err)
		}

		// push metrics
		countMetricsName := fmt.Sprintf("%s.count", monitor.MetricsPrefix)
		h.monitor.Incr(countMetricsName, []string{
			fmt.Sprintf("success:%v", isSuccess),
			fmt.Sprintf("path:%v", r.URL.Path),
			fmt.Sprintf("errorcode:%v", errorConverted.ECode),
		}, 1)

		gaugeMetricsName := fmt.Sprintf("%s.duration", monitor.MetricsPrefix)
		h.monitor.Gauge(gaugeMetricsName, float64(duration.Microseconds()), []string{
			fmt.Sprintf("path:%v", r.URL.Path),
		}, 1)

		// skip write response if response already written
		if httpCtx.ResponseAlreadyWrite {
			return
		}

		// default response
		if err != nil {
			httpCtx.WriteHTTPResponseErrorToJSON(err)
			return
		}

		httpCtx.WriteHTTPResponseToJSON(nil, http.StatusOK)
	}
}
