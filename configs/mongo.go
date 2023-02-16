package configs

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Database *mongo.Client = ConnectMongo()

func ConnectMongo() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(Env("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")

	coll := GetCollection(client, "promotions")

	model := mongo.IndexModel{Keys: bson.M{"shop": "text"}}

	name, err := coll.Indexes().CreateOne(ctx, model)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("MongoDB Indexed : ", name)

	return client
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("promotion").Collection(collectionName)
	return collection
}
