package opensearchclient

import (
	// golang package
	"crypto/tls"
	"net"
	"net/http"
	"time"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/secretmanager"

	// external package
	"github.com/opensearch-project/opensearch-go"
)

const (
	maxIdleConnectionPerHost       = 10
	responseHeaderTimeoutInSeconds = 5
)

// InitOpenSearchClient create opensearch client for future opensearch request
func InitOpenSearchClient(secrets secretmanager.SecretsOpenSearch) *opensearch.Client {
	cfg := opensearch.Config{
		Addresses: secrets.Addresses,
		Username:  secrets.Username,
		Password:  secrets.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   maxIdleConnectionPerHost,
			ResponseHeaderTimeout: time.Duration(responseHeaderTimeoutInSeconds) * time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Duration(responseHeaderTimeoutInSeconds) * time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS11,
				InsecureSkipVerify: false,
			},
		},
	}

	client, err := opensearch.NewClient(cfg)
	if err != nil {
		err = errors.Wrap(err).WithCode("OSC.IOSC00")
		log.Error(err, nil, "opensearch.NewClient() got error - InitOpenSearchClient")
	}

	return client
}
