package main

import (
	"database/sql"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"insights-api/handlers"
	"insights-api/services"
	"insights-api/tweets"
	"insights-api/user"
	"log"
	"net/http"
)

type Server struct {
	router        *mux.Router
	twitterClient *twitter.Client
	db            *sql.DB
}

func NewServer(router *mux.Router, twitterClient *twitter.Client, db *sql.DB) *Server {
	return &Server{router: router, twitterClient: twitterClient, db: db}
}

func (s Server) Start() error {
	s.registerRouteHandlers()

	log.Println("Server started on port 8000")

	handler := cors.Default().Handler(s.router)

	return http.ListenAndServe(":8000", handler)
}

func (s *Server) registerRouteHandlers() {
	serviceWrapper := services.NewServiceWrapper(s.twitterClient)
	timeline := tweets.NewTimelineService(serviceWrapper)
	users := user.NewUser(serviceWrapper)
	trackedUserDAO := user.NewTrackedUserDAO(s.db)

	s.router.HandleFunc("/user/{userId}/tweets", handlers.TweetsHandler(timeline)).Methods(http.MethodGet)
	s.router.HandleFunc("/users/query", handlers.UserSearchHandler(users)).Methods(http.MethodPost)
	s.router.HandleFunc("/users/tracking", handlers.GetTrackedUsersHandler(trackedUserDAO, users)).Methods(http.MethodGet)
	s.router.HandleFunc("/user/{userId}", handlers.UserLookupHandler(users)).Methods(http.MethodGet)
	s.router.HandleFunc("/user/track/new", handlers.TrackUserHandler(trackedUserDAO)).Methods(http.MethodPost)
	s.router.HandleFunc("/user/{userId}/untrack", handlers.UntrackUserHandler(trackedUserDAO)).Methods(http.MethodDelete)
}
