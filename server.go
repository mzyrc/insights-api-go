package main

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	handler := cors.Default().Handler(s.router)

	return http.ListenAndServe(":8000", handler)
}

func (s *Server) registerRouteHandlers() {
	serviceWrapper := services.NewServiceWrapper(s.twitterClient)
	timeline := services.NewTimelineService(serviceWrapper)
	users := services.NewUser(serviceWrapper)

	s.router.HandleFunc("/user/{userId}/tweets", handlers.TweetsHandler(timeline)).Methods(http.MethodGet)
	s.router.HandleFunc("/users/query", handlers.UserSearchHandler(users)).Methods(http.MethodPost)
	s.router.HandleFunc("/users/tracking", handlers.FollowedUsersHandler(users)).Methods(http.MethodGet)
	s.router.HandleFunc("/user/{userId}", handlers.UserLookupHandler(users)).Methods(http.MethodGet)
}
