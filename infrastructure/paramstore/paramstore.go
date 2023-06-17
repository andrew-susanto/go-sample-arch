package paramstore

import (
	// golang package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"

	// external package
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

//go:generate mockgen -source=paramstore.go -destination=paramstore_mock.go -package=paramstore

type ParamStoreClientInterface interface {
	// Get information about a single parameter by specifying the parameter name. To
	// get information about more than one parameter at a time, use the GetParameters
	// operation.
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

type ParamStore struct {
	client ParamStoreClientInterface
}

// InitParamstore inits aws paramstore client with given config
func InitParamstore(config aws.Config) ParamStore {
	client := ssm.NewFromConfig(config)

	return ParamStore{
		client: client,
	}
}

// GetValue gets value from aws paramstore with given key
//
// Returns value from aws paramstore if key exists
// Otherwise return empty string
func (ps ParamStore) GetValue(key string) string {
	input := &ssm.GetParameterInput{
		Name: aws.String(key),
	}

	result, err := ps.client.GetParameter(context.Background(), input)
	if err != nil {
		log.Error(err, key, "ps.client.GetParameter() got error - GetValue")
		return ""
	}

	if result.Parameter.Value == nil {
		log.Error(err, key, "result parameter value is nill - GetValue")
		return ""
	}

	return *result.Parameter.Value
}
