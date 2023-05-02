package main

import (
	"fmt"
	"net/http"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {

	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"To"`
		Subject string `json:"Subject"`
		Message string `json:"Message"`
	}
	var requestPayload mailMessage

	err := app.Tools.ReadJSON(w, r, &requestPayload)
	if err != nil {
		app.Tools.ErrorJSON(w, err)
		return
	}
	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendMTPMessage(msg)
	if err != nil {
		app.Tools.ErrorJSON(w, err)
		fmt.Println(err)
		return
	}

	payload := app.JSONResponse
	payload.Error = false
	payload.Message = "Sento to " + requestPayload.To

	app.Tools.WriteJSON(w, http.StatusAccepted, payload)

}
