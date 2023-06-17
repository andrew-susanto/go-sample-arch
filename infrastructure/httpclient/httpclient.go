package httpclient

import (
	// golang package
	"net/http"
	"time"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/config"

	// external package
	"github.com/afex/hystrix-go/hystrix"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const (
	maxHttpClientTimeoutInSeconds = 30

	defaultCircuitBreakerName                  = "default"
	defaultCircuitBreakeTimeoutInMicroSeconds  = 30000
	defaultCircuitBreakerConcurentRequests     = 100
	defaultCircuitBreakerErrorPercentThreshold = 25
)

// HttpClient is abstraction of http client
type HttpClient struct {
	config config.Config
	http   *http.Client
}

// InitHttpClient create http client for future http request
func InitHttpClient(config config.Config) HttpClient {
	client := &http.Client{
		Timeout:   time.Duration(maxHttpClientTimeoutInSeconds) * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	// register each circuit breaker config
	for circuitBreakerKey, circuitBreakerConfig := range config.HttpClient {
		hystrix.ConfigureCommand(circuitBreakerKey, hystrix.CommandConfig{
			Timeout:               circuitBreakerConfig.Timeout,
			MaxConcurrentRequests: circuitBreakerConfig.MaxConcurrentRequest,
			ErrorPercentThreshold: circuitBreakerConfig.ErrorPercentageThreshold,
		})
	}

	// use our own default circuit breaker rather than using hystrix default
	hystrix.ConfigureCommand(defaultCircuitBreakerName, hystrix.CommandConfig{
		Timeout:               defaultCircuitBreakeTimeoutInMicroSeconds,
		MaxConcurrentRequests: defaultCircuitBreakerConcurentRequests,
		ErrorPercentThreshold: defaultCircuitBreakerErrorPercentThreshold,
	})

	return HttpClient{
		config: config,
		http:   client,
	}
}

// Wrap will wraps function with given client name to enable circuit breaker
func (httpClient HttpClient) Wrap(clientName string, fn func() error) error {
	_, exists := httpClient.config.HttpClient[clientName]
	if !exists {
		clientName = defaultCircuitBreakerName
	}

	err := hystrix.Do(clientName, fn, nil)
	return err
}

// Do sends an HTTP request and returns an HTTP response, following policy (such as redirects, cookies, auth) as configured on the client.
func (httpClient HttpClient) Do(request *http.Request) (*http.Response, error) {
	return httpClient.http.Do(request)
}
