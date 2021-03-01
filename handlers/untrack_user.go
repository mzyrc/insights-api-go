package handlers

import (
	"fmt"
	"github.com/gorilla/mux"
	"insights-api/dao"
	"net/http"
	"strconv"
)

type UntrackUserDAO interface {
	Remove(userId int, twitterUserId int64) error
}

func UntrackUserHandler(untrackUserDAO UntrackUserDAO) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		currentUserId := 1

		requestVariables := mux.Vars(request)
		userId := requestVariables["userId"]

		twitterUserId, err := strconv.ParseInt(userId, 10, 64)

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Invalid user id in URL"))
			return
		}

		err = untrackUserDAO.Remove(currentUserId, twitterUserId)

		if err != nil {
			if userNotExistsErr, ok := err.(dao.UserNotExistsError); ok {
				writer.WriteHeader(http.StatusNotFound)
				writer.Write([]byte(fmt.Sprintf("User id %d is not tracked", userNotExistsErr.Id)))
				return
			} else {
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte("Unknown error occurred"))
				return
			}
		}

		writer.WriteHeader(http.StatusOK)
	}
}
