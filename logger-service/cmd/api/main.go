package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/cristianortiz/toolbox"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"   // to listen inciming request
	rpcPort  = "5001" // listen RPC calls
	mongoURL = "mongodb://mongo:27017"
	// mongoURL = "mongodb://localhost:27017" only to test if logger server works before  adding it to docker-compose
	grpcPort = "50001" //to listen gRPC calls

)

var client *mongo.Client

type Config struct {
	Models       data.Models
	Tools        toolbox.Tools
	JSONResponse *toolbox.JSONResponse
}

func main() {
	//connect to mongoDB
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}
	client = mongoClient

	//create a context in order to disconnect from mongo
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	//as usual when we are working with the context, defer cancel..
	defer cancel()
	//.. and then close the connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models:       data.New(client),
		JSONResponse: &toolbox.JSONResponse{},
	}
	//start RPC server for logger service
	//register the RPC server
	err = rpc.Register(new(RPCServer))
	go app.rpcListen()

	//start web server
	log.Println("Starting logger service in port", webPort)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// start to listen for rcp connection
func (app *Config) rpcListen() error {
	log.Println("Starting RPC server on port", rpcPort)
	//listen on RPC with standar library
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()
	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

// connectToMongo() create a mongoDB client options and  use it to create and return
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
		log.Println("Error connecting MongoDB", err)
		return nil, err
	}
	return conn, nil
}
