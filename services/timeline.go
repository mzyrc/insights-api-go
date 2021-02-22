package services

import (
	"errors"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"net/http"
)

type TimelineConfig struct {
	UserId int64
}

type TwitterTimelineClient interface {
	UserTimeline(config TimelineConfig) ([]twitter.Tweet, *http.Response, error)
}

func NewTimelineService(client TwitterTimelineClient) *Timeline {
	return &Timeline{client: client}
}

type Timeline struct {
	client TwitterTimelineClient
}

func (t Timeline) GetTimeLine(userId int64) ([]twitter.Tweet, error) {
	tweets, httpResponse, err := t.client.UserTimeline(TimelineConfig{UserId: userId})

	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
			return nil, errors.New("user does not exist")
		} else {
			return nil, errors.New(fmt.Sprintf("unexpected error occurred requesting timeline: %v", err.Error()))
		}
	}

	return tweets, nil
}
