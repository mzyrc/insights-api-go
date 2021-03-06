package tweets

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"net/http"
	"sort"
	"time"
)

type TimelineConfig struct {
	UserId int64
}

type twitterTimelineClient interface {
	UserTimeline(config TimelineConfig) ([]twitter.Tweet, *http.Response, error)
}

type tweetStorage interface {
	GetAll(userId int64) ([]LocalTweet, error)
	Save(tweets []LocalTweet) error
	SetLastSync(timestamp time.Time, tweet LocalTweet)
}

type tweetService struct {
	client       twitterTimelineClient
	tweetStorage tweetStorage
}

func NewTweetService(client twitterTimelineClient, db *sql.DB) *tweetService {
	dao := newTweetDAO(db)
	return &tweetService{client: client, tweetStorage: dao}
}

func (t tweetService) GetTimeLine(userId int64) ([]twitter.Tweet, error) {
	tweets, httpResponse, err := t.client.UserTimeline(TimelineConfig{UserId: userId})

	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
			return nil, errors.New("user does not exist")
		} else {
			return nil, errors.New(fmt.Sprintf("unexpected error occurred requesting tweetService: %v", err.Error()))
		}
	}

	return tweets, nil
}

func (t tweetService) GetTweets(userId int64) ([]LocalTweet, error) {
	tweets, err := t.tweetStorage.GetAll(userId)

	if err != nil {
		// @todo handle the error
	}

	return tweets, nil
}

func (t tweetService) StoreTweetsByUser(userId int64) error {
	tweetsByUser, err := t.GetTimeLine(userId)

	if err != nil {
		return err
	}

	tweets := make([]LocalTweet, len(tweetsByUser))

	for index, tweet := range tweetsByUser {
		tweets[index] = newTweetAdapter(&tweet).ToLocalTweet()
	}

	currentSyncTime := time.Now()
	err = t.tweetStorage.Save(tweets)

	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].ID > tweets[j].ID
	})

	if err != nil {
		return err
	}

	t.tweetStorage.SetLastSync(currentSyncTime, tweets[0])

	return nil
}
