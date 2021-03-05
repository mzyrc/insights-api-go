package user

import (
	"errors"
	"github.com/dghubble/go-twitter/twitter"
	"log"
	"net/http"
)

type twitterUserClient interface {
	Search(screenName string) ([]twitter.User, *http.Response, error)
	GetFriends(nextPageId int64) (*twitter.Friends, *http.Response, error)
	GetUser(userId int64) (*twitter.User, *http.Response, error)
	GetUsers(userIdList []int64) ([]twitter.User, *http.Response, error)
}

var ServiceUnavailable = errors.New("user service: Twitter API unavailable")

type User struct {
	client twitterUserClient
}

func NewUser(client twitterUserClient) *User {
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
		return nil, err
	}

	return twitterUsers, nil
}

func (u User) getFollowingList(data *[]twitter.User, nextPageId int64) error {
	result, httpResponse, err := u.client.GetFriends(nextPageId)

	if err != nil {
		if httpResponse != nil && (httpResponse.StatusCode == http.StatusTooManyRequests) {
			return ServiceUnavailable
		} else {
			return err
		}
	}

	*data = append(*data, result.Users...)

	if result.NextCursor == 0 {
		return nil
	}

	return u.getFollowingList(data, result.NextCursor)
}

func (u User) GetUser(userId int64) (*twitter.User, error) {
	user, _, err := u.client.GetUser(userId)

	if err != nil {
		// @todo do something
		log.Println(err.Error())
	}

	return user, nil
}

func (u User) GetUsers(userIdList []int64) ([]twitter.User, error) {
	users, _, err := u.client.GetUsers(userIdList)

	if err != nil {
		return nil, err
	}

	return users, nil
}
