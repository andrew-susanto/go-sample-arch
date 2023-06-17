package http

import (
	// golang package
	"net/http"
)

//go:generate mockgen -source=./repo.go -destination=./repo_mock.go -package=http

type HttpClient interface {
	// Wrap will wraps function with given client name to enable circuit breaker
	Wrap(clientName string, fn func() error) error

	// Do sends an HTTP request and returns an HTTP response, following policy (such as redirects, cookies, auth) as configured on the client.
	Do(req *http.Request) (*http.Response, error)
}

type Repository struct {
	httpClient HttpClient
}

// NewRepository is a function to initialize http repository
func NewRepository(client HttpClient) Repository {
	return Repository{
		httpClient: client,
	}
}
