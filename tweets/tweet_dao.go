package tweets

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
)

type Tweet struct {
	db *sql.DB
}

func NewTweetDAO(db *sql.DB) *Tweet {
	return &Tweet{db: db}
}

func (t Tweet) Save(tweets []LocalTweet) error {
	transaction, err := t.db.Begin()

	if err != nil {
		log.Println(err)
		return err
	}

	statement, statementErr := transaction.Prepare(
		pq.CopyIn("tweet", "tweet_id", "text", "user_id", "created_at", "favourite_count", "retweet_count"),
	)

	if statementErr != nil {
		log.Println(statementErr)
		return statementErr
	}

	for _, tweet := range tweets {
		_, err = statement.Exec(tweet, tweet.UserID, tweet.Text, tweet.CreatedAt, tweet.FavouriteCount, tweet.RetweetCount)

		if err != nil {
			log.Println(err)
			return err
		}
	}

	_, err = statement.Exec()

	if err != nil {
		log.Println(err)
		return err
	}

	err = statement.Close()

	if err != nil {
		log.Println(err)
		return err
	}

	err = transaction.Commit()

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
