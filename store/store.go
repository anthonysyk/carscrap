package store

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Collections struct {
	client *mongo.Client
	Cars   *mongo.Collection
}

func Init() *Collections {
	fmt.Println("Connecting to MongoDB")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://carscrap:mysecretpassword@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	cars := client.Database("carscrap").Collection("cars")

	return &Collections{
		client: client,
		Cars:   cars,
	}
}

func (c Collections) GetClient() *mongo.Client {
	return c.client
}
