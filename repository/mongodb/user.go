package mongodb

import (
	// golang package
	"context"
	"fmt"
	"time"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/errors"
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"

	// external package
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	sampleCollectionName  = "sample-collection"
	queryTimeoutInSeconds = 30
)

// InsertUser insert user to mongodb
//
// Returns inserted id and nil error if success
// Otherwise return empty string and non nil error
func (repo *Repository) InsertUser(ctx context.Context, user User) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutInSeconds*time.Second)
	defer cancel()

	collection := repo.mongo.Collection(sampleCollectionName)

	res, err := collection.InsertOne(ctx, bson.M{"name": "pi", "value": 3.14159})
	if err != nil {
		err = errors.Wrap(err).WithCode("RPST.GUBI00")
		log.Error(err, nil, "collection.InsertOne() got error - InsertUser")
		return "", err
	}

	insertedID := fmt.Sprintf("%v", res.InsertedID)
	return insertedID, nil
}

// GetUserByID gets user by given id
//
// Return user and nil error when success
// Otherwise return empty user and non nil error
func (repo *Repository) GetUserByID(ctx context.Context, id int64) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeoutInSeconds*time.Second)
	defer cancel()

	collection := repo.mongo.Collection(sampleCollectionName)

	cur, err := collection.Find(ctx, bson.D{{Key: "", Value: primitive.Regex{Pattern: "", Options: ""}}})
	if err != nil {
		err = errors.Wrap(err).WithCode("RPST.GUBI00")
		log.Error(err, nil, "collection.Find() got error - GetUserByID")
		return User{}, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			err = errors.Wrap(err).WithCode("RPST.GUBI01")
			log.Error(err, nil, "collection.Find() got error - GetUserByID")
			return User{}, err
		}
	}

	err = cur.Err()
	if err != nil {
		err = errors.Wrap(err).WithCode("RPST.GUBI02")
		log.Error(err, nil, "cur.Error() got error - GetUserByID")
		return User{}, err
	}

	return User{
		ID:     id,
		Name:   "dummy user name",
		Gender: 3,
	}, nil
}
