package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/mux"
	"insights-api/handlers"
	"insights-api/services"
	"log"
	"net/http"
)

type Server struct {
	router        *mux.Router
	twitterClient *twitter.Client
}

func NewServer(router *mux.Router, twitterClient *twitter.Client) *Server {
	return &Server{router: router, twitterClient: twitterClient}
}

func (s Server) Start() error {
	s.registerRouteHandlers()

	log.Println("Server started on port 8000")
	return http.ListenAndServe(":8000", s.router)
}

func (s *Server) registerRouteHandlers() {
	serviceWrapper := services.NewServiceWrapper(s.twitterClient)
	timeline := services.NewTimelineService(serviceWrapper)

	s.router.HandleFunc("/user/{userId}/tweets", handlers.TweetsHandler(timeline)).Methods(http.MethodGet)
}
