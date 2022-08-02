package helpers

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoOnce sync.Once

const (
	Db                 = "UserDb"
	UserCollectionName = "User"
)

var (
	UserCollection *mongo.Collection
)

type ConnectionHelper struct {
	connStr string
}

type MongoInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func New(connStr string) *ConnectionHelper {
	return &ConnectionHelper{
		connStr: connStr,
	}
}

func (c *ConnectionHelper) ConnectDb() {
	mongoOnce.Do(func() {

		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(c.connStr))

		if err != nil {
			panic(err)
		}

		err = client.Ping(context.TODO(), nil)

		if err != nil {
			panic(err)
		}

		db := client.Database(Db)

		UserCollection = db.Collection(UserCollectionName)
	})
}
