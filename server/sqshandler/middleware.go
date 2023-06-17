package sqshandler

import (
	// golang package
	"context"
	"fmt"
	"sync"
	"time"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/config"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/monitor"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/tracer"

	// external package
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSService interface {
	// Returns the URL of an existing Amazon SQS queue.
	GetQueueUrl(ctx context.Context, params *sqs.GetQueueUrlInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueUrlOutput, error)

	// Retrieves one or more messages (up to 10)
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)

	// Deletes the specified message from the specified queue.
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

// registerQueueConsumer registers queue consumer based on given config
func (h *Handler) registerQueueConsumer(ctx context.Context, wg *sync.WaitGroup, monitorService monitor.Monitor, sqsConfig config.SQSClientConfig, sqsService SQSService, fn func(context.Context, types.Message) error) {
	if !sqsConfig.Enabled {
		return
	}

	urlResult, err := sqsService.GetQueueUrl(context.Background(), &sqs.GetQueueUrlInput{
		QueueName: &sqsConfig.QueueName,
	})
	if err != nil {
		err = errors.Wrap(err).WithCode("MDL.HF00")
		log.Error(err, sqsConfig, "sqsService.GetQueueUrl() got error - handleFunc")
		return
	}

	// prevent blocking
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msgResult, err := sqsService.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
					MessageAttributeNames: []string{string(types.QueueAttributeNameAll)},
					QueueUrl:              urlResult.QueueUrl,
					MaxNumberOfMessages:   int32(sqsConfig.MaxNumberMessage),
				})
				if err != nil {
					err = errors.Wrap(err).WithCode("MDL.HF01")
					log.Error(err, sqsConfig, "sqsService.ReceiveMessage() got error - handleFunc")
				}

				// process message in async if no error
				if err == nil {
					wg.Add(1)
					go func() {
						defer wg.Done()
						for _, message := range msgResult.Messages {
							// recover from panic to prevent server crash
							defer func() {
								r := recover()
								if r != nil {
									switch t := r.(type) {
									case error:
										err = errors.Wrap(err).WithCode("MDL.HF02")
									default:
										err = errors.New(fmt.Sprintf("%v", t)).WithCode("MDL.HF03")
									}
								}
							}()

							ctx, span := tracer.Start(context.Background(), sqsConfig.QueueName)
							defer span.End()

							start := time.Now()
							err := fn(ctx, message)
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
								errorConverted = errors.Wrap(errConvert).WithCode("MDL.HF04")
							}

							// push metrics
							countMetricsName := fmt.Sprintf("%s.count", monitor.MetricsPrefix)
							monitorService.Incr(countMetricsName, []string{
								fmt.Sprintf("success:%v", isSuccess),
								fmt.Sprintf("queuename:%v", sqsConfig.QueueName),
								fmt.Sprintf("errorcode:%v", errorConverted.ECode),
							}, 1)

							gaugeMetricsName := fmt.Sprintf("%s.duration", monitor.MetricsPrefix)
							monitorService.Gauge(gaugeMetricsName, float64(duration.Microseconds()), []string{
								fmt.Sprintf("queuename:%v", sqsConfig.QueueName),
							}, 1)

							if !isSuccess {
								continue
							}

							// delete message from queue if success
							_, err = sqsService.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
								QueueUrl:      urlResult.QueueUrl,
								ReceiptHandle: message.ReceiptHandle,
							})
							if err != nil {
								err = errors.Wrap(err).WithCode("MDL.HF05")
								log.Error(err, message, "sqsService.DeleteMessageWithContext() got error - registerQueueConsumer")
							}
						}
					}()
				}
			}

			time.Sleep(time.Duration(sqsConfig.PollPeriodInMilisecond) * time.Millisecond)
		}
	}()
}
