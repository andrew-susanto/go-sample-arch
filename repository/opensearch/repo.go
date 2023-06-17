package opensearch

import "net/http"

//go:generate mockgen -source=./repo.go -destination=./repo_mock.go -package=opensearch

type OpenSearchClient interface {
	Perform(*http.Request) (*http.Response, error)
}

type Repository struct {
	opensearch OpenSearchClient
}

// NewRepository is a function to initialize opensearch repository
func NewRepository(opensearch OpenSearchClient) Repository {
	return Repository{
		opensearch: opensearch,
	}
}
