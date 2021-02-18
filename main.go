package main

import (
	"errors"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	godotenv.Load()

	twitterClient, err := NewTwitterClient()

	if err != nil {
		log.Fatal("Could not initialise a Twitter client")
	}

	server := NewServer(mux.NewRouter(), twitterClient)
	log.Fatal(server.Start())
}

func NewTwitterClient() (*twitter.Client, error) {
	if os.Getenv("TWITTER_API_KEY") == "" {
		return nil, errors.New("missing TWITTER_API_KEY env variable")
	}

	if os.Getenv("TWITTER_API_SECRET") == "" {
		return nil, errors.New("missing TWITTER_API_SECRET env variable")
	}

	if os.Getenv("TWITTER_API_ACCESS_TOKEN") == "" {
		return nil, errors.New("missing TWITTER_API_ACCESS_TOKEN env variable")
	}

	if os.Getenv("TWITTER_API_ACCESS_TOKEN_SECRET") == "" {
		return nil, errors.New("missing TWITTER_API_ACCESS_TOKEN_SECRET env variable")
	}

	config := oauth1.NewConfig(os.Getenv("TWITTER_API_KEY"), os.Getenv("TWITTER_API_SECRET"))
	token := oauth1.NewToken(os.Getenv("TWITTER_API_ACCESS_TOKEN"), os.Getenv("TWITTER_API_ACCESS_TOKEN_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient), nil
}
