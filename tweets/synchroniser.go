package tweets

import (
	"fmt"
	"log"
	"time"
)

type tweetSyncStorage interface {
	Save(tweets []Tweet) error
	GetDueSyncList() ([]SyncConfig, error)
	UpdateLastSync(twitterUserId int64, lastTweetId int64, lastSync time.Time) error
}

type tweetSyncType interface {
	StoreTweetsByUser(userId int64, lastTweetId int64) ([]Tweet, error)
}

type synchroniser struct {
	tweetDAO    tweetSyncStorage
	tweetSync   tweetSyncType
}

func newSynchroniser(dao tweetSyncStorage, tweetSync tweetSyncType) *synchroniser {
	return &synchroniser{
		tweetDAO:  dao,
		tweetSync: tweetSync,
	}
}

func (s synchroniser) Start() {
	log.Println("Starting Synchroniser")
	ticker := time.NewTicker(time.Minute * 10)

	s.synchronise()

	for {
		select {
		case <-ticker.C:
			s.synchronise()
		}
	}
}

func (s synchroniser) synchronise() {
	log.Println("Starting synchronisation cycle")
	syncList, err := s.tweetDAO.GetDueSyncList()

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(fmt.Sprintf("Found %d users to synchronise tweets", len(syncList)))

	for _, sync := range syncList {
		lastTweetId := sync.LastTweetId
		tweets, storeTweetsErr := s.tweetSync.StoreTweetsByUser(sync.UserID, lastTweetId)

		if storeTweetsErr != nil {
			log.Println(err)
		}

		if len(tweets) > 0 {
			lastTweetId = tweets[0].ID
		}

		s.tweetDAO.UpdateLastSync(sync.UserID, lastTweetId, time.Now())
	}

	log.Println("Finished synchronisation cycle")
}
