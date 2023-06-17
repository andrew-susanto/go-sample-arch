package cronhandler

import (
	// golang package
	"context"
	"fmt"
	"sync"
	"time"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/monitor"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/tracer"
)

// registerCron registers cron function with given monnitor, cron name, cron interval in seconds, and cron function
func (h *Handler) registerCron(ctx context.Context, wg *sync.WaitGroup, monitorService monitor.Monitor, cronName string, cronIntervalInSeconds int, fn func(context.Context) error) {
	// prevent blocking
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				wg.Add(1)
				// run cron function in async
				go func() {
					defer wg.Done()
					var err error

					// recover from panic to prevent server crash
					defer func() {
						r := recover()
						if r != nil {
							switch t := r.(type) {
							case error:
								err = errors.Wrap(err).WithCode("MDL.HF00")
							default:
								err = errors.New(fmt.Sprintf("%v", t)).WithCode("MDL.HF01")
							}
						}
					}()

					// replace context with otel tracer context
					spanName := fmt.Sprintf("%s - interval %d", cronName, cronIntervalInSeconds)
					ctx, span := tracer.Start(context.Background(), spanName)
					defer span.End()

					start := time.Now()
					err = fn(ctx)
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
						errorConverted = errors.Wrap(errConvert).WithCode("MDL.HF02")
					}

					// push metrics
					countMetricsName := fmt.Sprintf("%s.count", monitor.MetricsPrefix)
					monitorService.Incr(countMetricsName, []string{
						fmt.Sprintf("success:%v", isSuccess),
						fmt.Sprintf("cronname:%v", cronName),
						fmt.Sprintf("errorcode:%v", errorConverted.ECode),
					}, 1)

					gaugeMetricsName := fmt.Sprintf("%s.duration", monitor.MetricsPrefix)
					monitorService.Gauge(gaugeMetricsName, float64(duration.Microseconds()), []string{
						fmt.Sprintf("cronname:%v", cronName),
					}, 1)
				}()
			}

			time.Sleep(time.Duration(cronIntervalInSeconds) * time.Second)
		}
	}()
}
