package main

import (
	"net/http"
	"toolbox"
)

//Broker handle function
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	var tools toolbox.Tools

	payload := toolbox.JSONResponse{
		Error:   false,
		Message: "hiting the  Broker..",
	}
	_ = tools.WriteJSON(w, http.StatusOK, payload)
}
