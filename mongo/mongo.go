package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	client *mongo.Client
	db     *mongo.Database
}

func New(dbName string, opt *options.ClientOptions) *Mongo {
	c, e := mongo.Connect(context.Background(), opt)
	if e != nil {
		panic(e.Error())
	}
	
	n := new(Mongo)
	n.client = c
	n.db = n.client.Database(dbName)
	
	return n
}

func (m *Mongo) CloseConnection() error {
	return m.client.Disconnect(context.Background())
}

func (m *Mongo) CollectionList(filter bson.M) []string {
	l, e := m.db.ListCollectionNames(context.Background(), filter, options.ListCollections().
		SetNameOnly(true).
		SetAuthorizedCollections(true))
	
	if e != nil {
		return nil
	}
	return l
}

func (m *Mongo) InsertOne(collection string, data interface{}) (*mongo.InsertOneResult, error) {
	return m.db.Collection(collection).InsertOne(context.Background(), data)
}

func (m *Mongo) InsertMany(collection string, data []interface{}) (*mongo.InsertManyResult, error) {
	return m.db.Collection(collection).InsertMany(context.Background(), data)
}

func (m *Mongo) CreateCollection(key string, opts ...*options.CreateCollectionOptions) error {
	return m.db.CreateCollection(context.Background(), key, opts...)
}

func (m *Mongo) CreateCollectionIndexMany(collection string, indexes []mongo.IndexModel, opts ...*options.CreateIndexesOptions) ([]string, error) {
	return m.db.Collection(collection).Indexes().CreateMany(context.Background(), indexes, opts...)
}

func (m *Mongo) CreateCollectionIndexOne(collection string, index mongo.IndexModel, opts ...*options.CreateIndexesOptions) (string, error) {
	return m.db.Collection(collection).Indexes().CreateOne(context.Background(), index, opts...)
}
