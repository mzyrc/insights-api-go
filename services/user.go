package services

import (
	"errors"
	"github.com/dghubble/go-twitter/twitter"
	"log"
	"net/http"
)

type TwitterUserClient interface {
	Search(screenName string) ([]twitter.User, *http.Response, error)
	GetFriends(nextPageId int64) (*twitter.Friends, *http.Response, error)
}

type User struct {
	client TwitterUserClient
}

func NewUser(client TwitterUserClient) *User {
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

func (u User) GetFollowingList() ([]twitter.User, error) {
	var twitterUsers []twitter.User
	err := u.getFollowingList(&twitterUsers, 0)

	if err != nil {
		return nil, errors.New("what do we do here?")
	}

	return twitterUsers, nil
}

func (u User) getFollowingList(data *[]twitter.User, nextPageId int64) error {
	result, httpResponse, err := u.client.GetFriends(nextPageId)

	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusInternalServerError {

		}
	}

	*data = append(*data, result.Users...)

	if result.NextCursor == 0 {
		return nil
	}

	return u.getFollowingList(data, result.NextCursor)
}
