package docdb

import (
	// golang package
	"context"
	"fmt"
	"time"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/secretmanager"

	// external package
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Timeout operations after N seconds
	connectTimeout = 5
	queryTimeout   = 30

	// Which instances to read from
	readPreference = "secondaryPreferred"

	connectionStringTemplate = "mongodb://%s:%s@%s/%s?replicaSet=rs0&readpreference=%s"
)

// InitDocDB inits mongodb client based on given secrets
func InitDocDB(secrets secretmanager.SecretsDocDb) *mongo.Client {
	connectionURI := fmt.Sprintf(connectionStringTemplate, secrets.Username, secrets.Password, secrets.ClusterEndpoint, secrets.DBName, readPreference)

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionURI))
	if err != nil {
		err = errors.Wrap(err).WithCode("DCDB.IDB01")
		log.Fatal(err, nil, "mongo.NewClient() got error - InitDocDB")
		return client
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)
	defer cancel()

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		err = errors.Wrap(err).WithCode("DCDB.IDB03")
		log.Error(err, nil, "client.Ping() got error - InitDocDB")
		return client
	}

	log.Info(nil, "success connect to mongodb database - InitDocDB")
	return client
}
