package main

import (
	"context"
	"log"
	"logger-service/data"
	"time"
)

// RPC specific type
type RPCServer struct{}

// payload to receive from RPC
type RPCPayload struct {
	Name string
	Data string
}

// LogInfo writes the payload to mongoDB
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {

	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})

	if err != nil {
		log.Println("error writing to mongo", err)
		return err
	}

	*resp = "Processed payload via RPC " + payload.Name
	return nil

}
