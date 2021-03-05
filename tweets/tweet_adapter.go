package tweets

import "github.com/dghubble/go-twitter/twitter"

type LocalTweet struct {
	ID             int64  `json:"id"`
	CreatedAt      string `json:"created_at"`
	Text           string `json:"text"`
	UserID         int64  `json:"user_id"`
	FavouriteCount int    `json:"favourite_count"`
	RetweetCount   int    `json:"retweet_count"`
}

type tweetAdapter struct {
	tweet *twitter.Tweet
}

func newTweetAdapter(tweet *twitter.Tweet) *tweetAdapter {
	return &tweetAdapter{tweet: tweet}
}

func (t tweetAdapter) ToLocalTweet() LocalTweet {
	return LocalTweet{
		ID:             t.tweet.ID,
		CreatedAt:      t.tweet.CreatedAt,
		Text:           t.tweet.Text,
		UserID:         t.tweet.User.ID,
		FavouriteCount: t.tweet.FavoriteCount,
		RetweetCount:   t.tweet.RetweetCount,
	}
}
