package tweets

import (
	"database/sql"
	"github.com/dghubble/go-twitter/twitter"
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
	CalculateSentiment(userId int64)
}

type TweetService struct {
	tweetStorage tweetStorage
	Synchroniser *synchroniser
	TweetSync    *tweetSync
}

func NewTweetService(client twitterTimelineClient, db *sql.DB) *TweetService {
	dao := newTweetDAO(db)
	tweetHTTPClient := newTweetClient(client)
	tweetSynchroniser := newTweetSync(tweetHTTPClient, dao)

	return &TweetService{
		tweetStorage: dao,
		Synchroniser: newSynchroniser(dao, tweetSynchroniser),
		TweetSync:    tweetSynchroniser,
	}
}

func (t TweetService) GetTweets(userId int64) ([]Tweet, error) {
	tweets, err := t.tweetStorage.GetAll(userId)

	if err != nil {
		// @todo handle the error
	}

	return tweets, nil
}

func (t TweetService) StoreTweetsForFirstTime(userId int64) {
	now := time.Now()
	tweets, err := t.TweetSync.StoreTweetsByUser(userId, 0)

	if err != nil {
		//	@todo do something useful
	}

	t.tweetStorage.SetLastSync(now, tweets[0])
}
