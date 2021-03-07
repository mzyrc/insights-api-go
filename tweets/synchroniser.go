package tweets

import (
	"fmt"
	"log"
	"sort"
	"time"
)

type tweetSyncStorage interface {
	Save(tweets []Tweet) error
	GetDueSyncList() ([]SyncConfig, error)
	UpdateLastSync(twitterUserId int64, lastTweetId int64, lastSync time.Time) error
}

type synchroniser struct {
	tweetClient TweetHTTPClient
	tweetDAO    tweetSyncStorage
}

func newSynchroniser(service TweetHTTPClient, dao tweetSyncStorage) *synchroniser {
	return &synchroniser{
		tweetClient: service,
		tweetDAO:    dao,
	}
}

func (s synchroniser) Start() {
	log.Println("Starting Synchroniser")
	ticker := time.NewTicker(time.Minute * 10)

	s.begin()

	for {
		select {
		case <-ticker.C:
			s.begin()
		}
	}
}

func (s synchroniser) begin() {
	log.Println("Starting synchronisation cycle")
	syncList, err := s.tweetDAO.GetDueSyncList()

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(fmt.Sprintf("Found %d users to synchronise tweets", len(syncList)))

	for _, sync := range syncList {
		lastTweetId := sync.LastTweetId
		tweets, err := s.tweetClient.GetTimeLine(TimelineConfig{UserId: sync.UserID, StartFromTweetID: sync.LastTweetId})

		if err != nil {
			log.Println(err)
			return
		}

		tweetsByUser := make([]Tweet, len(tweets))

		for index, tweet := range tweets {
			tweetsByUser[index] = newTweetAdapter(&tweet).ToLocalTweet()
		}

		lastSync := time.Now()
		err = s.tweetDAO.Save(tweetsByUser)

		if err != nil {
			log.Println(err)
		}

		if len(tweetsByUser) > 0 {
			sort.Slice(tweetsByUser, func(i, j int) bool {
				return tweetsByUser[i].ID > tweetsByUser[j].ID
			})
			lastTweetId = tweetsByUser[0].ID
		}

		err = s.tweetDAO.UpdateLastSync(sync.UserID, lastTweetId, lastSync)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println("Finished synchronisation cycle")
}
