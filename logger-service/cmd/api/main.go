package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"   // to listen
	rpcPort  = "5001" // listen RPC calls
	mongoURL = "mongodb://mongo:27017"
	grpcPort = "50001" //to listen gRPC calls

)

var client *mongo.Client

type Config struct {
}

func main() {
	//connect to mongoDB
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	//create a context in order to disconnect from DB
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//as usual when we are working with the context, defer cancel..
	defer cancel()
	//.. and then close the connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

//connectToMongo() create a mongoDB client options and the use it to create and return
// a mongoDB connect client type
func connectToMongo() (*mongo.Client, error) {
	//create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	//connecto to DB
	conn, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting DB", err)
		return nil, err
	}
	return conn, nil
}
