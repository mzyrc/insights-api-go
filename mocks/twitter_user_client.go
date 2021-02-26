package mocks

import (
	"github.com/dghubble/go-twitter/twitter"
	"net/http"
)

type MockTwitterUserClient struct {
	MockSearch     func(screenName string) ([]twitter.User, *http.Response, error)
	MockGetFriends func(nextPageId int64) (*twitter.Friends, *http.Response, error)
	MockGetUser    func(userId int64) (*twitter.User, *http.Response, error)
}

func (m MockTwitterUserClient) Search(screenName string) ([]twitter.User, *http.Response, error) {
	return m.MockSearch(screenName)
}

func (m MockTwitterUserClient) GetFriends(nextPageId int64) (*twitter.Friends, *http.Response, error) {
	return m.MockGetFriends(nextPageId)
}

func (m MockTwitterUserClient) GetUser(userId int64) (*twitter.User, *http.Response, error) {
	return m.MockGetUser(userId)
}
