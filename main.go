package main

import (
	"context"
	"log"
	"os"

	"github.com/BerkCicekler/e-commerce-audio-api/cmd/api"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClient *mongo.Client

func init() {
	//load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("env load error", err)
	}
	
	log.Println("env file loaded")

	// create mongodb client
	mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))


	if err != nil {
		log.Fatal("connection error", err)
	}

	err = mongoClient.Ping(context.Background(), readpref.Nearest());
	if err != nil {
		log.Fatal("ping error", err)
	}

	log.Println("mongodb connected")
}

func main() {
	defer mongoClient.Disconnect(context.Background());

	mongoDb := mongoClient.Database(os.Getenv("DB_NAME"))

	api := api.NewAPIServer(":8080")
	api.Run(mongoDb)
}