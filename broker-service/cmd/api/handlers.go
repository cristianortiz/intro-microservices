package main

import (
	"net/http"
)

//Broker handle function
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "hiting the  Broker..",
	}
	_ = app.writeJSON(w, http.StatusOK, payload)

}
