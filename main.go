package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	godotenv.Load()

	twitterClient, err := NewTwitterClient()

	if err != nil {
		log.Fatal(err.Error())
	}

	db, dbErr := NewDBClient()

	if dbErr != nil {
		log.Fatal(dbErr.Error())
	}

	server := NewServer(mux.NewRouter(), twitterClient, db)
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

func NewDBClient() (*sql.DB, error) {
	if os.Getenv("DB_USER") == "" {
		return nil, errors.New("missing DB_USER env variable")
	}

	user := os.Getenv("DB_USER")

	if os.Getenv("DB_PASSWORD") == "" {
		return nil, errors.New("missing DB_PASSWORD env variable")
	}

	password := os.Getenv("DB_PASSWORD")

	if os.Getenv("DB_HOST") == "" {
		return nil, errors.New("missing DB_HOST env variable")
	}

	host := os.Getenv("DB_HOST")

	if os.Getenv("DB_PORT") == "" {
		return nil, errors.New("missing DB_PORT env variable")
	}

	port := os.Getenv("DB_PORT")

	if os.Getenv("DB_NAME") == "" {
		return nil, errors.New("missing DB_NAME env variable")
	}

	dbName := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		return nil, err
	}

	_, err = db.Query("SELECT 1+1")

	if err != nil {
		return nil, err
	}

	return db, nil
}
