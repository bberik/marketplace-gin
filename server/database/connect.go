package database

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func NewMongoDBClient() (client *mongo.Client) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	mongo_url := os.Getenv("MONGODB_URL")

	client, err = mongo.NewClient(options.Client().ApplyURI(mongo_url))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("MongoDB client created successfully")
	return
}

var Client *mongo.Client = NewMongoDBClient()

func NewCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)

	return collection
}

var Users *mongo.Collection = NewCollection(Client, "users")
var Products *mongo.Collection = NewCollection(Client, "products")
var Orders *mongo.Collection = NewCollection(Client, "orders")
var Shops *mongo.Collection = NewCollection(Client, "shops")
