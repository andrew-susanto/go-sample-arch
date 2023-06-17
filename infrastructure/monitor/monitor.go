package monitor

import (
	// golang package
	"fmt"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"

	// external package
	"github.com/DataDog/datadog-go/v5/statsd"
)

const (
	datadogHost   = "127.0.0.1:8125"
	MetricsPrefix = "cxp.crm_app"
)

type Monitor interface {
	// Incr is just Count of 1
	Incr(name string, tags []string, rate float64) error

	// Gauge measures the value of a metric at a particular time.
	Gauge(name string, value float64, tags []string, rate float64) error
}

// InitMonitor inits metrics client
func InitMonitor(environment string, serviceName string) Monitor {
	// universal service tags
	tagsOption := statsd.WithTags([]string{
		fmt.Sprintf("env:%v", environment),
		fmt.Sprintf("service:%v", serviceName),
	})

	dogstatsd_client, err := statsd.New(datadogHost, tagsOption)
	if err != nil {
		log.Fatal(err, map[string]interface{}{
			"datadogHost": datadogHost,
		}, "statsd.New() got error - InitMonitor")
	}

	return dogstatsd_client
}
