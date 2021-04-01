package tweets

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"log"
	"net/http"
)

type twitterTimelineClient interface {
	UserTimeline(config TimelineConfig) ([]twitter.Tweet, *http.Response, error)
}

type tweetClient struct {
	client twitterTimelineClient
}

func newTweetClient(client twitterTimelineClient) *tweetClient {
	return &tweetClient{client: client}
}

func (t tweetClient) GetTimeLine(config TimelineConfig) ([]twitter.Tweet, error) {
	if config.UserId == 0 {
		return nil, errors.New("must specify a a UserID property in the config")
	}

	log.Println(fmt.Sprintf("Fetching tweets for user id: %d", config.UserId))
	tweets, httpResponse, err := t.client.UserTimeline(config)

	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
			return nil, errors.New("user does not exist")
		} else {
			return nil, errors.New(fmt.Sprintf("unexpected error occurred requesting TweetService: %v", err.Error()))
		}
	}

	return tweets, nil
}

func (t tweetClient) CalculateSentiment(userId int64) {
	jsonBody, _ := json.Marshal(map[string]int64{"user_id": userId})

	response, err := http.Post("http://localhost:8001/user/sentiment/calculate", "application/json", bytes.NewBuffer(jsonBody))

	if err != nil {
		log.Println(err)
	}

	defer response.Body.Close()
}
