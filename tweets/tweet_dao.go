package tweets

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"log"
	"time"
)

type SyncConfig struct {
	LastTweetId    int64     `json:"last_tweet_id"`
	UserID         int64     `json:"user_id"`
	SynchronisedAt time.Time `json:"synchronised_at"`
}

type tweetDAO struct {
	db *sql.DB
}

func newTweetDAO(db *sql.DB) *tweetDAO {
	return &tweetDAO{db: db}
}

func (t tweetDAO) Save(tweets []Tweet) error {
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

	log.Println(fmt.Sprintf("Successfully stored %d tweets", len(tweets)))

	return nil
}

func (t tweetDAO) GetAll(userId int64) ([]Tweet, error) {
	log.Println("Fetching tweets from the database")
	rows, err := t.db.Query("SELECT * FROM tweet WHERE user_id = $1", userId)

	if err != nil {
		return nil, err
	}

	var tweets []Tweet

	for rows.Next() {
		var tweet Tweet
		err = rows.Scan(&tweet.ID, &tweet.Text, &tweet.UserID, &tweet.CreatedAt, &tweet.FavouriteCount, &tweet.RetweetCount)
		if err != nil {
			return nil, err
		}

		tweets = append(tweets, tweet)
	}

	return tweets, nil
}

func (t tweetDAO) SetLastSync(timestamp time.Time, tweet Tweet) {
	insertSQL := "INSERT INTO tweet_synchronisation (last_tweet_id, user_id, synchronised_at) VALUES ($1, $2, $3)"
	_, err := t.db.Exec(insertSQL, tweet.ID, tweet.UserID, timestamp)

	if err != nil {
		log.Println(err)
	}

	log.Println("Set the last sync")
}

func (t tweetDAO) GetDueSyncList() ([]SyncConfig, error) {
	rows, err := t.db.Query("SELECT last_tweet_id, user_id, synchronised_at FROM tweet_synchronisation WHERE synchronised_at <= NOW() - INTERVAL '10 minutes'")

	if err != nil {
		return nil, err
	}

	var syncList []SyncConfig

	for rows.Next() {
		var sync SyncConfig
		err = rows.Scan(&sync.LastTweetId, &sync.UserID, &sync.SynchronisedAt)

		if err != nil {
			log.Println(err)
		}

		syncList = append(syncList, sync)
	}

	return syncList, nil
}

func (t tweetDAO) UpdateLastSync(twitterUserId int64, lastTweetId int64, lastSync time.Time) error {
	updateSQL := "UPDATE tweet_synchronisation SET last_tweet_id = $1, synchronised_at = $2 WHERE user_id = $3"
	_, err := t.db.Exec(updateSQL, lastTweetId, lastSync, twitterUserId)

	if err != nil {
		return err
	}

	return nil
}
