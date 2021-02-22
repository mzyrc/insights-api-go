package handlers

import (
	"encoding/json"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/mux"
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
			sendErrorResponse(writer, http.StatusBadRequest, "invalid user id in URL")
			return
		}

		tweets, timelineErr := service.GetTimeLine(twitterUserId)

		if timelineErr != nil {
			sendErrorResponse(writer, http.StatusBadRequest, timelineErr.Error())
			return
		}

		sendResponseData(writer, tweets)
	}
}

func sendErrorResponse(writer http.ResponseWriter, statusCode int, errorMessage string) {
	writer.WriteHeader(statusCode)
	writer.Write([]byte(errorMessage))
	return
}

func sendResponseData(writer http.ResponseWriter, tweets []twitter.Tweet) {
	response, jsonErr := json.Marshal(tweets)

	if jsonErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(jsonErr.Error()))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}
