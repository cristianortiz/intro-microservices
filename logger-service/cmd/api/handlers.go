package main

import (
	"logger-service/data"
	"net/http"
)

type JSONPayload struct {
	Name string `string:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {

	//read json into var
	var requestPayload JSONPayload
	err := app.Tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.Tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	//insert data
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		app.Tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	//send back
	payload := app.JSONResponse
	payload.Error = false
	payload.Message = "logged!"

	app.Tools.WriteJSON(w, http.StatusAccepted, payload)
}
