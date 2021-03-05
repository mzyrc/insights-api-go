package handlers

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/mux"
	tweets2 "insights-api/tweets"
	"net/http"
	"strconv"
)

type TimelineService interface {
	GetTimeLine(userId int64) ([]twitter.Tweet, error)
}

func TweetsHandler(service TimelineService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestVariables := mux.Vars(request)
		userId := requestVariables["userId"]

		twitterUserId, err := strconv.ParseInt(userId, 10, 64)

		if err != nil {
			respondWithError(writer, http.StatusBadRequest, "invalid user id in URL")
			return
		}

		tweets, timelineErr := service.GetTimeLine(twitterUserId)
		userTweets := make([]tweets2.LocalTweet, len(tweets))

		for index, tweet := range tweets {
			userTweets[index] = tweets2.NewTweet(&tweet).ToLocalTweet()
		}

		if timelineErr != nil {
			respondWithError(writer, http.StatusBadRequest, timelineErr.Error())
			return
		}

		respondWithSuccess(writer, http.StatusOK, userTweets)
	}
}
