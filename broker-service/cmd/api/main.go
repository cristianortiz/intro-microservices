package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cristianortiz/toolbox"
)

const webPort = "80"

type Config struct {
	Tools        toolbox.Tools
	JSONResponse *toolbox.JSONResponse
}

func main() {
	var tools toolbox.Tools
	app := Config{
		Tools:        tools,
		JSONResponse: &toolbox.JSONResponse{},
	}

	log.Printf(("Starting broker service in port: %s"), webPort)

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
