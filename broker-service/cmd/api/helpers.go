package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// struct to send back response as a JSON
type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	//thanks to go 1.18 generics
	Data any `json:"da       ta,omitempty"`
}

// readJson receiver function to Config type to reading JSON
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	// one megabyte limitation for upload json file
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}
	//check for only a single json value in the received file
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}
	return nil
}

// writeJson receiver function to Config type to write a  JSON file, the header arg is optional
func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, header ...http.Header) error {

	out := json.NewEncoder(w)
	//check if header was included as the last parameter
	if len(header) > 0 {
		for key, value := range header[0] {
			w.Header()[key] = value

		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := out.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}
	payload := jsonResponse{
		Error:   true,
		Message: err.Error(),
	}
	return app.writeJSON(w, statusCode, payload)
}
