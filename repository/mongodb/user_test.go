package mongodb

import (
	// golang package
	"context"
	"testing"

	// external package
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestRepository_GetUserByID(t *testing.T) {
	type args struct {
		userID int64
	}

	tests := []struct {
		name     string
		args     args
		want     User
		wantErr  error
		mockFunc func(*mtest.T)
	}{
		{
			name: "when_get_user_by_id_success_then_xxx",
			args: args{
				userID: 123,
			},
			want: User{
				ID:     123,
				Name:   "dummy user name",
				Gender: 3,
			},
			wantErr: nil,
			mockFunc: func(mt *mtest.T) {
				first := mtest.CreateCursorResponse(1, "dbname.sample-collection", mtest.FirstBatch, bson.D{{Key: "test", Value: "test"}})
				getMore := mtest.CreateCursorResponse(1, "dbname.sample-collection", mtest.NextBatch, bson.D{{Key: "test", Value: "test"}})
				killCursors := mtest.CreateCursorResponse(0, "dbname.sample-collection", mtest.NextBatch)
				mt.AddMockResponses(first, getMore, killCursors)
			},
		},
	}

	for _, tt := range tests {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		defer mt.Close()

		mt.Run(tt.name, func(mt *mtest.T) {
			// log.InitLogger()
			repository := NewRepository(mt.Client.Database("dbname"))
			tt.mockFunc(mt)

			got, gotErr := repository.GetUserByID(context.Background(), tt.args.userID)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, gotErr)
		})
	}
}
