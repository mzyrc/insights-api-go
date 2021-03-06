package services

import (
	"github.com/dghubble/go-twitter/twitter"
	"insights-api/tweets"
	"net/http"
)

type ServiceWrapper struct {
	client *twitter.Client
}

func NewTweetServiceWrapper(client *twitter.Client) *ServiceWrapper {
	return &ServiceWrapper{client: client}
}

func (s ServiceWrapper) UserTimeline(params tweets.TimelineConfig) ([]twitter.Tweet, *http.Response, error) {
	return s.client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID: params.UserId,
		//ScreenName:      "",
		//Count:           0,
		//SinceID:         0,
		//MaxID:           0,
		//TrimUser:        nil,
		//ExcludeReplies:  nil,
		//IncludeRetweets: nil,
		//TweetMode:       "",
	})
}

func (s ServiceWrapper) Search(screenName string) ([]twitter.User, *http.Response, error) {
	return s.client.Users.Search(screenName, &twitter.UserSearchParams{
		Query: screenName,
		//Page:            0,
		//Count:           0,
		//IncludeEntities: nil,
	})
}

func (s ServiceWrapper) GetFriends(nextPageId int64) (*twitter.Friends, *http.Response, error) {
	return s.client.Friends.List(&twitter.FriendListParams{
		UserID: 0,
		//ScreenName:          "",
		Cursor: nextPageId,
		//Count:               0,
		//SkipStatus:          nil,
		//IncludeUserEntities: nil,
	})
}

func (s ServiceWrapper) GetUser(userId int64) (*twitter.User, *http.Response, error) {
	return s.client.Users.Show(&twitter.UserShowParams{
		UserID: userId,
		//ScreenName:      "",
		//IncludeEntities: nil,
	})
}

func (s ServiceWrapper) GetUsers(userIdList []int64) ([]twitter.User, *http.Response, error) {
	return s.client.Users.Lookup(&twitter.UserLookupParams{
		UserID: userIdList,
		//ScreenName:      nil,
		//IncludeEntities: nil,
	})
}
