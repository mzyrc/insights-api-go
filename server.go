package main

import (
	"database/sql"
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
	twitterClient *services.ServiceWrapper
	db            *sql.DB
	tweetService  *tweets.TweetService
}

func NewServer(router *mux.Router, twitterClient *services.ServiceWrapper, db *sql.DB, tweetService *tweets.TweetService) *Server {
	return &Server{
		router:        router,
		twitterClient: twitterClient,
		db:            db,
		tweetService:  tweetService,
	}
}

func (s Server) Start() error {
	s.registerRouteHandlers()

	log.Println("Server started on port 8000")

	handler := cors.Default().Handler(s.router)

	return http.ListenAndServe(":8000", handler)
}

func (s *Server) registerRouteHandlers() {
	tweetService := tweets.NewTweetService(s.twitterClient, s.db)
	users := user.NewUser(s.twitterClient)
	trackedUserDAO := user.NewTrackedUserDAO(s.db)

	s.router.HandleFunc("/user/{userId}/tweets", handlers.TweetsHandler(tweetService)).Methods(http.MethodGet)
	s.router.HandleFunc("/users/query", handlers.UserSearchHandler(users)).Methods(http.MethodPost)
	s.router.HandleFunc("/users/tracking", handlers.GetTrackedUsersHandler(trackedUserDAO, users)).Methods(http.MethodGet)
	s.router.HandleFunc("/user/{userId}", handlers.UserLookupHandler(users)).Methods(http.MethodGet)
	s.router.HandleFunc("/user/track/new", handlers.TrackUserHandler(trackedUserDAO, tweetService)).Methods(http.MethodPost)
	s.router.HandleFunc("/user/{userId}/untrack", handlers.UntrackUserHandler(trackedUserDAO)).Methods(http.MethodDelete)
}
