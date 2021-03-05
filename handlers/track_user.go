package handlers

import (
	"encoding/json"
	"insights-api/user"
	"net/http"
	"strconv"
)

type trackUserPostBody struct {
	UserId string `json:"user_id"`
}

type trackedUserDAO interface {
	Create(userId int, twitterUserId int64) error
}

type tweetSyncService interface {
	StoreTweetsByUser(userId int64) error
}

func TrackUserHandler(trackedUser trackedUserDAO, tweetSync tweetSyncService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		currentUser := 1

		var postBody trackUserPostBody
		requestBodyDecoder := json.NewDecoder(request.Body)
		decoderErr := requestBodyDecoder.Decode(&postBody)

		if decoderErr != nil {
			respondWithError(writer, http.StatusBadRequest, "Invalid post body received")
			return
		}

		twitterUserId, err := strconv.ParseInt(postBody.UserId, 10, 64)

		if err != nil {
			respondWithError(writer, http.StatusBadRequest, "Invalid user_id in post body")
			return
		}

		err = trackedUser.Create(currentUser, twitterUserId)

		if err != nil {
			if userExistsErr, ok := err.(user.UserExistsError); ok {
				respondWithError(writer, http.StatusConflict, userExistsErr.Error())
				return
			} else {
				respondWithError(writer, http.StatusInternalServerError, "An unknown error occurred")
				return
			}
		}

		go tweetSync.StoreTweetsByUser(twitterUserId)

		respondWithSuccess(writer, http.StatusOK, "")
	}
}
