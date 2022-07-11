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

//auth type to handle the payload request from auth microservice
type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Broker handle function to check if broker service is up
func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	payload := toolbox.JSONResponse{
		Error:   false,
		Message: "hiting the  Broker..",
	}
	_ = app.Tools.WriteJSON(w, http.StatusOK, payload)
}

//Handlesubmission() handler function to handle the request of any microservice
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	//to put the json converted body request
	var requestPayload RequestPayload
	//read the data in request body and convert into json
	err := app.Tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.Tools.ErrorJSON(w, err)
	}
	//check wich microservice is requesting the broker
	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		app.Tools.ErrorJSON(w, errors.New("unkown action"))
	}
}

//authenticate() call auth microservice method with their own payload
func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	//create some JSON w'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//call the microservice
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
	//check if the response json is setting with errors
	if jsonFromService.Error {
		app.Tools.ErrorJSON(w, err, http.StatusUnauthorized)
		return
	}
	//send back to the user the auth payload
	payload := app.JSONResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.Tools.WriteJSON(w, http.StatusAccepted, payload)

}
