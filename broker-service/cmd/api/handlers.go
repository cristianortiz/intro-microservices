package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"toolbox"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Broker handle function
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := toolbox.JSONResponse{
		Error:   false,
		Message: "hiting the  Broker..",
	}
	_ = app.Tools.WriteJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.Tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.Tools.ErrorJSON(w, err)
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		app.Tools.ErrorJSON(w, errors.New("unkown action"))
	}

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	//create some JSON w'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.Tools.ErrorJSON(w, err)
		return
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.Tools.ErrorJSON(w, err)
		return
	}
	defer response.Body.Close()

	//make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.Tools.ErrorJSON(w, errors.New("invalid credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.Tools.ErrorJSON(w, errors.New("error calling auth service"))
		return
	}

	//create a variable we'll read response.body into
	jsonFromService := app.JSONResponse

	//decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.Tools.ErrorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		app.Tools.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := app.JSONResponse
	payload.Error = false
	payload.Message = "Authenticaded!"
	payload.Data = jsonFromService.Data

	app.Tools.WriteJSON(w, http.StatusAccepted, payload)

}
