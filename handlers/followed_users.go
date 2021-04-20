package handlers

import (
	"errors"
	"github.com/dghubble/go-twitter/twitter"
	"insights-api/user"
	"net/http"
)

type trackedUsersDAO interface {
	GetUsers(userId int) ([]int64, error)
}

type userLookUpService interface {
	GetUsers(userIdList []int64) ([]twitter.User, error)
}

var UnexpectedError = errors.New("unexpected cant find user error")

func GetTrackedUsersHandler(trackedUserDAO trackedUsersDAO, service TwitterUserService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		currentUserId := 1
		userIdList, err := trackedUserDAO.GetUsers(currentUserId)

		if err != nil {
			// @todo do something useful here
		}

		if len(userIdList) == 0 {
			respondWithSuccess(writer, http.StatusOK, []user.LocalUser{})
			return
		}

		trackedUsers, err := getTrackedUsers(service, userIdList)

		respondWithSuccess(writer, http.StatusOK, trackedUsers)
	}
}

func getTrackedUsers(service userLookUpService, userIdList []int64) ([]user.LocalUser, error) {
	twitterUsers, err := service.GetUsers(userIdList)

	if err != nil {
		// @todo do something useful
		return nil, UnexpectedError
	}

	trackedUsers := make([]user.LocalUser, len(twitterUsers))

	for index, twitterUser := range twitterUsers {
		trackedUsers[index] = user.NewUserAdapter(twitterUser).ToLocalUser()
	}
	return trackedUsers, nil
}
