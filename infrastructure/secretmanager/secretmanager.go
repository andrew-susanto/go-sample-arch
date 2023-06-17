package secretmanager

import (
	// golang package
	"context"
	"encoding/json"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"

	// external package
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var (
	secretsName = "cxp-crm.secrets"
)

// InitSecretManager initialize aws secrets manager
func InitSecretManager(config aws.Config) Secrets {
	client := secretsmanager.NewFromConfig(config)
	result, err := client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
		SecretId: &secretsName,
	})
	if err != nil {
		log.Fatal(err, nil, "client.GetSecretValue() got error - InitSecretmanager")
		return Secrets{}
	}

	var secret Secrets
	err = json.Unmarshal([]byte(*result.SecretString), &secret)
	if err != nil {
		log.Fatal(err, nil, "json.Umarhsal() got error - InitSecretManager")
		return Secrets{}
	}

	return secret
}
