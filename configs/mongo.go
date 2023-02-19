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

	coll_promotions := GetCollection(client, "promotions")

	model_promotions := mongo.IndexModel{Keys: bson.M{"shop": "text"}}

	coll_admin := GetCollection(client, "admins")

	model_admins := mongo.IndexModel{Keys: bson.M{"email": -1}, Options: options.Index().SetUnique(true)}

	index_promotions, index_promotions_error := coll_promotions.Indexes().CreateOne(ctx, model_promotions)
	if index_promotions_error != nil {
		log.Fatal(err)
	}

	index_admins, index_admins_error := coll_admin.Indexes().CreateOne(ctx, model_admins)
	if index_admins_error != nil {
		log.Fatal(err)
	}

	log.Println("MongoDB Indexed : ", index_promotions)
	log.Println("MongoDB Indexed : ", index_admins)

	return client
}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("promotion").Collection(collectionName)
	return collection
}
