package opensearchclient

import (
	// golang package
	"testing"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/secretmanager"

	// external package
	"github.com/stretchr/testify/assert"
)

func TestOpenSearchClient_InitOpenSearch(t *testing.T) {
	secrets := secretmanager.SecretsOpenSearch{
		Addresses: []string{"test"},
		Domain:    "test",
		Username:  "test",
		Password:  "test",
	}

	got := InitOpenSearchClient(secrets)
	assert.NotNil(t, got)
}
