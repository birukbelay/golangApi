package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoCollection struct {

}

func MongodbConn() (client *mongo.Client) {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		panic(err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB Connected successfully !")
	return client

}

//func init() {
//	mongoDbConn := MongodbConn()
//	defer CloseMongo(mongoDbConn)
//
//
//}

//var ItemCollection = mongoDbConn.Database("database1").Collection("items")



func CloseMongo(client *mongo.Client) {
	err := client.Disconnect(context.TODO())

	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Connection to MongoDB closed.")
}


