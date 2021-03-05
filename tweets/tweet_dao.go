package tweets

import (
	"database/sql"
	"github.com/lib/pq"
	"log"
)

type tweetDAO struct {
	db *sql.DB
}

func newTweetDAO(db *sql.DB) *tweetDAO {
	return &tweetDAO{db: db}
}

func (t tweetDAO) Save(tweets []LocalTweet) error {
	transaction, err := t.db.Begin()

	if err != nil {
		log.Println(err)
		return err
	}

	statement, statementErr := transaction.Prepare(
		pq.CopyIn("tweet", "id", "text", "user_id", "created_at", "favourite_count", "retweet_count"),
	)

	if statementErr != nil {
		log.Println(statementErr)
		return statementErr
	}

	for _, tweet := range tweets {
		_, err = statement.Exec(tweet.ID, tweet.Text, tweet.UserID, tweet.CreatedAt, tweet.FavouriteCount, tweet.RetweetCount)

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

	log.Println("Successfully stored tweets")

	return nil
}

func (t tweetDAO) GetAll(userId int64) ([]LocalTweet, error) {
	log.Println("Fetching tweets from the database")
	rows, err := t.db.Query("SELECT * FROM tweet WHERE user_id = $1", userId)

	if err != nil {
		return nil, err
	}

	var tweets []LocalTweet

	for rows.Next() {
		var tweet LocalTweet
		err = rows.Scan(&tweet.ID, &tweet.Text, &tweet.UserID, &tweet.CreatedAt, &tweet.FavouriteCount, &tweet.RetweetCount)
		if err != nil {
			return nil, err
		}

		tweets = append(tweets, tweet)
	}

	return tweets, nil
}
