package tweets

import (
	"database/sql"
	"github.com/dghubble/go-twitter/twitter"
	"sort"
	"time"
)

type tweetStorage interface {
	GetAll(userId int64) ([]Tweet, error)
	Save(tweets []Tweet) error
	SetLastSync(timestamp time.Time, tweet Tweet)
}

type TimelineConfig struct {
	UserId           int64
	StartFromTweetID int64
}

type TweetHTTPClient interface {
	GetTimeLine(config TimelineConfig) ([]twitter.Tweet, error)
}

type TweetService struct {
	client       twitterTimelineClient
	tweetStorage tweetStorage
	tweetClient  TweetHTTPClient
	Synchroniser *synchroniser
}

func NewTweetService(client twitterTimelineClient, db *sql.DB) *TweetService {
	dao := newTweetDAO(db)
	tweetClient := newTweetClient(client)

	return &TweetService{
		client:       client,
		tweetStorage: dao,
		tweetClient:  tweetClient,
		Synchroniser: newSynchroniser(tweetClient, dao),
	}
}

func (t TweetService) GetTweets(userId int64) ([]Tweet, error) {
	tweets, err := t.tweetStorage.GetAll(userId)

	if err != nil {
		// @todo handle the error
	}

	return tweets, nil
}

func (t TweetService) StoreTweetsByUser(userId int64) error {
	tweetsByUser, err := t.tweetClient.GetTimeLine(TimelineConfig{UserId: userId})

	if err != nil {
		return err
	}

	tweets := make([]Tweet, len(tweetsByUser))

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
