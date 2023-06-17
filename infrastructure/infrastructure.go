package infrastructure

import (
	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/config"
)

//go:generate mockgen -source=infrastructure.go -destination=infrastructure_mock.go -package=infrastructure

// ParamStore is an interface for getting value from paramstore
type ParamStore interface {
	// GetValue is a function to get value from paramstore
	GetValue(key string) string
}

// Infrastructure is a struct that contains all infrastructure
type Infrastructure struct {
	Config     config.Config
	paramstore ParamStore
}

// InitInfrastructure is a function to initialize infrastructure
func InitInfrastructure(config config.Config, paramstore ParamStore) Infrastructure {
	return Infrastructure{
		Config:     config,
		paramstore: paramstore,
	}
}

// GetParamStoreValue is a function to get value from paramstore
func (infra *Infrastructure) GetParamStoreValue(key string) string {
	return infra.paramstore.GetValue(key)
}
