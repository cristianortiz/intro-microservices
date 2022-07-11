package main

import (
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
	payload := app.JSONResponse
	payload.Error = false
	payload.Message = fmt.Sprintf("Logged in user %s", user.Email)
	payload.Data = user

	app.Tools.WriteJSON(w, http.StatusAccepted, payload)
}
