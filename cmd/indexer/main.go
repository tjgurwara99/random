package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	dbName           string
	index            string
	unique           bool
	listOnly         bool
	collectionName   string
	connectionString string
)

func connectClient(connectionString string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client, err
}

func createIndex(dbName, collectionName, index string, unique bool) {
	client, err := connectClient(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()
	collection := client.Database(dbName).Collection(collectionName)
	_, err = collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{index: 1},
		Options: options.Index().SetUnique(unique),
	})

	if err != nil {
		log.Fatal(fmt.Errorf("Error occured while creating index: %w", err))
	}
}

func listIndices(dbName, collectionName string) {
	client, err := connectClient(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()
	collection := client.Database(dbName).Collection(collectionName)
	indexes, err := collection.Indexes().List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	var results []bson.M
	if err = indexes.All(ctx, &results); err != nil {
		log.Fatal(err)
	}
	for _, res := range results {
		fmt.Printf("%v\n", res)
	}
}

func init() {
	flag.StringVar(&connectionString, "addr", "mongodb://localhost:27017", "mongodb connection string")
	flag.StringVar(&dbName, "db", "test", "database name")
	flag.StringVar(&collectionName, "collection", "test", "collection name")
	flag.StringVar(&index, "index", "name", "index name for creating a new index")
	flag.BoolVar(&unique, "unique", false, "unique index")
	flag.BoolVar(&listOnly, "list-only", true, "only list all indices without changing")

}

func main() {
	flag.Parse()
	if !listOnly {
		createIndex(dbName, collectionName, index, unique)
	}
	listIndices(dbName, collectionName)
}
