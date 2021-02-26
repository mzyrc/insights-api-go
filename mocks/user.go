package mocks

import "github.com/dghubble/go-twitter/twitter"

type MockUser struct {
	MockSearch           func(screenName string) (users []twitter.User, err error)
	MockGetFollowingList func() ([]twitter.User, error)
	MockGetUser func(userId int64) (*twitter.User, error)
}

func (m MockUser) Search(screenName string) (users []twitter.User, err error) {
	return m.MockSearch(screenName)
}

func (m MockUser) GetFollowingList() ([]twitter.User, error) {
	return m.MockGetFollowingList()
}

func (m MockUser) GetUser(userId int64) (*twitter.User, error) {
	return m.MockGetUser(userId)
}