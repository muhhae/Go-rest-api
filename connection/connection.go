package connection

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var user_col *mongo.Collection
func User() *mongo.Collection {
	return user_col
}

var post_col *mongo.Collection
func Post() *mongo.Collection {
	return post_col
}

func InitConnection() {
	mongo_host := os.Getenv("MONGO_HOST")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongo_host))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	user_col = client.Database("golang").Collection("user")
	post_col = client.Database("golang").Collection("post")

	log.Println("Connected to MongoDB")
}
