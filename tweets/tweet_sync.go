package tweets

import (
	"sort"
)

type tweetSync struct {
	tweetClient  TweetHTTPClient
	tweetStorage tweetStorage
}

func newTweetSync(tweetClient TweetHTTPClient, tweetStorage tweetStorage) *tweetSync {
	return &tweetSync{
		tweetClient:  tweetClient,
		tweetStorage: tweetStorage,
	}
}

func (t tweetSync) StoreTweetsByUser(userId int64, lastTweetId int64) ([]Tweet, error) {
	config := TimelineConfig{UserId: userId}

	if lastTweetId != 0 {
		config.StartFromTweetID = lastTweetId
	}

	tweetsByUser, err := t.tweetClient.GetTimeLine(config)

	if err != nil {
		return nil, err
	}

	tweets := make([]Tweet, len(tweetsByUser))

	for index, tweet := range tweetsByUser {
		tweets[index] = newTweetAdapter(&tweet).ToLocalTweet()
	}

	err = t.tweetStorage.Save(tweets)

	sort.Slice(tweets, func(i, j int) bool {
		return tweets[i].ID > tweets[j].ID
	})

	if err != nil {
		return nil, err
	}

	t.tweetClient.CalculateSentiment(userId)

	return tweets, nil
}
