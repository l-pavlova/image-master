package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// This is a method to close resources.
// This method closes mongoDB connection and cancel context.
func Close(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

// This is a method that returns mongo.Client,
// context.Context, context.CancelFunc and error.
func Connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {

	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

// InsertOne is a method, used to insert
func InsertOne(client *mongo.Client, ctx context.Context, dataBase, col string, doc interface{}) (*mongo.InsertOneResult, error) {

	collection := client.Database(dataBase).Collection(col)

	result, err := collection.InsertOne(ctx, doc)
	return result, err
}

// InsertMany and error if any.
func InsertMany(client *mongo.Client, ctx context.Context, dataBase, col string, docs []interface{}) (*mongo.InsertManyResult, error) {

	collection := client.Database(dataBase).Collection(col)

	// InsertMany accepts two argument of type Context
	// and of empty interface
	result, err := collection.InsertMany(ctx, docs)
	return result, err
}

// query method returns a cursor and error.
func Query(client *mongo.Client, ctx context.Context,
	dataBase, col string, query, field interface{}) (result *mongo.Cursor, err error) {

	collection := client.Database(dataBase).Collection(col)

	result, err = collection.Find(ctx, query,
		options.Find().SetProjection(field))
	return
}

// find a single element matching a query
func FindOne(client *mongo.Client, ctx context.Context, dataBase, col string, query, field interface{}) (result *mongo.SingleResult) {
	collection := client.Database(dataBase).Collection(col)

	result = collection.FindOne(ctx, query, options.FindOne().SetProjection(field))
	return
}

// A single document that match with the
// filter will get updated.
// update contains the filed which should get updated.
func UpdateOne(client *mongo.Client, ctx context.Context, dataBase,
	col string, filter, update interface{}) (result *mongo.UpdateResult, err error) {

	// select the database and the collection
	collection := client.Database(dataBase).Collection(col)

	result, err = collection.UpdateOne(ctx, filter, update)
	return
}

// deleteOne is a user defined function that deletes
// a single document from the collection.
func DeleteOne(client *mongo.Client, ctx context.Context,
	dataBase, col string, query interface{}) (result *mongo.DeleteResult, err error) {

	// select document and collection
	collection := client.Database(dataBase).Collection(col)

	// query is used to match a document  from the collection.
	result, err = collection.DeleteOne(ctx, query)
	return
}

// deleteMany is a user defined function that deletes
// multiple documents from the collection.
func DeleteMany(client *mongo.Client, ctx context.Context,
	dataBase, col string, query interface{}) (result *mongo.DeleteResult, err error) {

	// select document and collection
	collection := client.Database(dataBase).Collection(col)

	// query is used to match  documents  from the collection.
	result, err = collection.DeleteMany(ctx, query)
	return
}
