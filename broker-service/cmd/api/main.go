package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/cristianortiz/toolbox"
	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Tools        toolbox.Tools         //access for toolbox methods
	JSONResponse *toolbox.JSONResponse //toolbox struct
	Rabbit       *amqp.Connection
}

func main() {
	//connect to rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()
	var tools toolbox.Tools
	app := Config{
		Tools:        tools,
		JSONResponse: &toolbox.JSONResponse{},
		Rabbit:       rabbitConn,
	}

	log.Printf(("Starting broker service in port: %s"), webPort)

	//define the http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	//start the server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	//set a back off logic in case of slow start of rabbitMQ
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbitMQ is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready..")
			counts++
		} else {
			connection = c
			log.Println("Connected to RabbitMQ..!")

			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}
		//if no error is returned but several retrys to connect has been executed
		if counts > 5 && err == nil {
			log.Printf("retried more than %v times", counts)
			return nil, fmt.Errorf("retried more than %v times", counts)
		}
		//doubling the back off time in every retry to connect
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off..")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
