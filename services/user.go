package services

import (
	"github.com/dghubble/go-twitter/twitter"
	"log"
	"net/http"
)

type UserProcessor interface {
	Search(screenName string) ([]twitter.User, *http.Response, error)
}

type User struct {
	client UserProcessor
}

func NewUser(client UserProcessor) *User {
	return &User{client: client}
}

func (u User) Search(screenName string) ([]twitter.User, error) {
	users, httpResponse, err := u.client.Search(screenName)

	if err != nil {
		log.Println(httpResponse)
		// @todo do something useful here
	}
	return users, nil
}
