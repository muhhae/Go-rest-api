package connection

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var user_col *mongo.Collection

func User() *mongo.Collection {
	return user_col
}

func InitConnection() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	log.Println("Connecting to MongoDB, result:", client)
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	user_col = client.Database("golang").Collection("user")
	log.Println("Connected to MongoDB")
}
