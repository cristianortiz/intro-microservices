package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	//using chi as a router
	mux := chi.NewRouter()

	//specify who is allowed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	//heart beat to check if service is still alive
	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/", app.Broker)
	//Boker single endpoint to handle request for any other microservices
	mux.Post("/handle", app.HandleSubmission)
	return mux
}
