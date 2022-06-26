package main

import "net/http"

//this is an alternative version of helper just for ilustration purposes
//to send back response as a JSON
type jsonResponseTest struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	//using interface rater than 1.18 generics 'any' type
	Data interface{} `json:"data,omitempty"`
}

// readJson receiver function to Config type to reading JSON
func (app *Config) readJSONTest(w http.ResponseWriter, r *http.Request, data interface{}) {

}
