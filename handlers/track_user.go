package handlers

import (
	"encoding/json"
	"insights-api/dao"
	"net/http"
	"strconv"
)

type trackUserPostBody struct {
	UserId string `json:"user_id"`
}

type trackedUserDAO interface {
	Create(userId int, twitterUserId int64) error
}

func TrackUserHandler(trackedUser trackedUserDAO) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		currentUser := 1

		var postBody trackUserPostBody
		requestBodyDecoder := json.NewDecoder(request.Body)
		decoderErr := requestBodyDecoder.Decode(&postBody)

		if decoderErr != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Invalid post body received"))
			return
		}

		twitterUserId, err := strconv.ParseInt(postBody.UserId, 10, 64)

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Invalid user_id in post body"))
			return
		}

		err = trackedUser.Create(currentUser, twitterUserId)

		if err != nil {
			if userExistsErr, ok := err.(dao.UserExistsError); ok {
				writer.WriteHeader(http.StatusConflict)
				writer.Write([]byte(userExistsErr.Error()))
				return
			} else {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte("An unknown error occurred"))
				return
			}
		}

		writer.WriteHeader(http.StatusOK)
	}
}
