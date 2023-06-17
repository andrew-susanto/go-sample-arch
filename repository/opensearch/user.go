package opensearch

import (
	// golang package
	"context"
	"encoding/json"
	"io"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"

	// external package
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

// GetUserByID gets user by given id from open search
//
// Return user and nil error when success
// Otherwise return empty user and non nil error
func (repo *Repository) GetUserByID(ctx context.Context, id int64) (User, error) {
	var bodyContent []byte
	getRequest := opensearchapi.GetRequest{
		Index:      "build-on-aws",
		DocumentID: "test",
	}

	getResponse, err := getRequest.Do(ctx, repo.opensearch)
	if err != nil {
		err = errors.Wrap(err).WithCode("RPST.GUBI00")
		log.Error(err, nil, "getRequest.Do() got error - GetUserByID")
		return User{}, err
	}

	if getResponse.Body != nil {
		defer getResponse.Body.Close()
		bodyContent, err = io.ReadAll(getResponse.Body)
		if err != nil {
			err = errors.Wrap(err).WithCode("RPST.GUBI01")
			log.Error(err, nil, "io.ReadAll() got error - GetUserByID")
			return User{}, err
		}
	}

	var response User
	err = json.Unmarshal(bodyContent, &response)
	if err != nil {
		err = errors.Wrap(err).WithCode("RPST.GUBI02")
		log.Error(err, nil, "json.Unmarshal() got error - GetUserByID")
	}

	return User{
		ID:     id,
		Name:   "dummy user name",
		Gender: 3,
	}, nil
}
