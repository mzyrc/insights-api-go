package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"insights-api/user"
	"net/http"
	"strconv"
)

type untrackUserDAO interface {
	Remove(userId int, twitterUserId int64) error
}

func UntrackUserHandler(untrackUserDAO untrackUserDAO) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		currentUserId := 1

		requestVariables := mux.Vars(request)
		userId := requestVariables["userId"]

		twitterUserId, err := strconv.ParseInt(userId, 10, 64)

		if err != nil {
			respondWithError(writer, http.StatusBadRequest, "Invalid user id in URL")
			return
		}

		err = untrackUserDAO.Remove(currentUserId, twitterUserId)

		if err != nil {
			if userNotExistsErr, ok := err.(user.UserNotExistsError); ok {
				respondWithError(writer, http.StatusNotFound, fmt.Sprintf("User id %d is not tracked", userNotExistsErr.Id))
				return
			} else {
				respondWithError(writer, http.StatusInternalServerError, "Unknown error occurred")
				return
			}
		}

		respondWithSuccess(writer, http.StatusOK, "")
	}
}
