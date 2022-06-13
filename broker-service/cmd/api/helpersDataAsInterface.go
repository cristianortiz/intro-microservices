package main

import "net/http"

//to send back response as a JSON
type jsonResponseTest struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	//thanks to go 1.18 generics
	Data interface{} `json:"data,omitempty"`
}

// readJson receiver function to Config type to reading JSON
func (app *Config) readJSONTest(w http.ResponseWriter, r *http.Request, data interface{}) {

}
