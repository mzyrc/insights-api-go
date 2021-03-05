package handlers

import (
	"github.com/gorilla/mux"
	"insights-api/tweets"
	"net/http"
	"strconv"
)

type tweetService interface {
	GetTweets(userId int64) ([]tweets.LocalTweet, error)
}

func TweetsHandler(service tweetService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestVariables := mux.Vars(request)
		userId := requestVariables["userId"]

		twitterUserId, err := strconv.ParseInt(userId, 10, 64)

		if err != nil {
			respondWithError(writer, http.StatusBadRequest, "invalid user id in URL")
			return
		}

		userTweets, timelineErr := service.GetTweets(twitterUserId)

		if timelineErr != nil {
			respondWithError(writer, http.StatusBadRequest, timelineErr.Error())
			return
		}

		respondWithSuccess(writer, http.StatusOK, userTweets)
	}
}
