package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.Tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.Tools.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	//rquest user data  and validate  againts DB
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		app.Tools.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	//validate password but send back same error msg if is not valid
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.Tools.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	//log authentication
	err = app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.Tools.ErrorJSON(w, err)
		return
	}
	payload := app.JSONResponse
	payload.Error = false
	payload.Message = fmt.Sprintf("Logged in user %s", user.Email)
	payload.Data = user

	app.Tools.WriteJSON(w, http.StatusAccepted, payload)
}

// logRequest sends a request to logger service to log the user succesfully authentication
func (app *Config) logRequest(name, data string) error {

	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
