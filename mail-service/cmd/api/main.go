package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/cristianortiz/toolbox"
)

type Config struct {
	Mailer       Mail
	Tools        toolbox.Tools
	JSONResponse *toolbox.JSONResponse
}

const webPort = "80"

func main() {

	app := Config{
		Mailer:       createMail(),
		JSONResponse: &toolbox.JSONResponse{},
	}

	log.Println("Starting mail service on port", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

func createMail() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	m := Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromName:    os.Getenv("FROM_NAME"),
		FromAddress: os.Getenv("FROM_ADDRESS"),
	}

	return m
}
