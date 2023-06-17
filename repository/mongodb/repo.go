package mongodb

import (
	// external package
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//go:generate mockgen -source=./repo.go -destination=./repo_mock.go -package=mongodb

type MongoDB interface {
	// Collection gets a handle for a collection with the given name configured with the given CollectionOptions.
	Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
}

type Repository struct {
	mongo MongoDB
}

// NewRepository is a function to initialize mongodb repository
func NewRepository(mongo MongoDB) Repository {
	return Repository{
		mongo: mongo,
	}
}
