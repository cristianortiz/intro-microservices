package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cristianortiz/toolbox"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int64

type Config struct {
	DB           *sql.DB
	Models       data.Models
	Tools        toolbox.Tools
	JSONResponse *toolbox.JSONResponse
}

func main() {
	log.Println("Starting authentication service..")

	// connect to DB
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres..")
	}
	// setup config
	app := Config{
		DB:           conn,
		Models:       data.New(conn),
		JSONResponse: &toolbox.JSONResponse{},
	}
	//define the http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	//start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// openDB take a string defined connection and open a  connection to postgres DB,
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// connectToDB checks if the docker postgres container is up and running before open the DB connection
func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	//infinite loop, keep going until the postgres container is connected
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgress is no yet ready..")
			counts++
		} else {
			log.Println("Connected to PostgresDB")
			return connection
		}
		//if failed attempts to connect to DB is more than 10, stop for two secs and retry
		if counts > 10 {
			log.Println(err)
			return nil
		}
		log.Println("Backing off for two secods...")
		time.Sleep(2 * time.Second)
		//just for readability, not mandatory
		continue
	}

}
