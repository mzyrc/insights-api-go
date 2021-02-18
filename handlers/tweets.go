package handlers

import (
	"encoding/json"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type TimelineProcessor interface {
	GetTimeLine(userId int64) ([]twitter.Tweet, error)
}

func TweetsHandler(service TimelineProcessor) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestVariables := mux.Vars(request)
		userId := requestVariables["userId"]

		twitterUserId, err := strconv.ParseInt(userId, 10, 64)

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode([]byte("invalid user id in URL"))
			return
		}

		tweets, timelineErr := service.GetTimeLine(twitterUserId)

		if timelineErr != nil {
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode([]byte(timelineErr.Error()))
			return
		}

		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(tweets)
	}
}
