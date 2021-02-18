package services

import (
	"github.com/dghubble/go-twitter/twitter"
	"net/http"
)

type ServiceWrapper struct {
	client *twitter.Client
}

func NewServiceWrapper(client *twitter.Client) *ServiceWrapper {
	return &ServiceWrapper{client: client}
}

func (s ServiceWrapper) UserTimeline(params TimelineConfig) ([]twitter.Tweet, *http.Response, error) {
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
