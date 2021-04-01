package tweets

import (
	"database/sql"
	"encoding/json"
	"github.com/dghubble/go-twitter/twitter"
)

type NullableFloat64 struct {
	sql.NullFloat64
}

func (nf *NullableFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

type Tweet struct {
	ID             int64           `json:"id"`
	CreatedAt      string          `json:"created_at"`
	Text           string          `json:"text"`
	UserID         int64           `json:"user_id"`
	FavouriteCount int             `json:"favourite_count"`
	RetweetCount   int             `json:"retweet_count"`
	SentimentScore NullableFloat64 `json:"sentiment_score"`
}

type tweetAdapter struct {
	tweet *twitter.Tweet
}

func newTweetAdapter(tweet *twitter.Tweet) *tweetAdapter {
	return &tweetAdapter{tweet: tweet}
}

func (t tweetAdapter) ToLocalTweet() Tweet {
	return Tweet{
		ID:             t.tweet.ID,
		CreatedAt:      t.tweet.CreatedAt,
		Text:           t.tweet.Text,
		UserID:         t.tweet.User.ID,
		FavouriteCount: t.tweet.FavoriteCount,
		RetweetCount:   t.tweet.RetweetCount,
	}
}
