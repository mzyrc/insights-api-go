package handlers

import (
	"insights-api/adapters"
	"net/http"
)

type GetTrackedUsersDAO interface {
	GetUsers(userId int) ([]int64, error)
}

func GetTrackedUsersHandler(trackedUserDAO GetTrackedUsersDAO, service TwitterUserService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		currentUserId := 1
		userIdList, err := trackedUserDAO.GetUsers(currentUserId)

		if err != nil {
			// @todo do something useful here
		}

		if len(userIdList) == 0 {
			respondWithSuccess(writer, http.StatusOK, []adapters.LocalUser{})
			return
		}

		twitterUsers, err := service.GetUsers(userIdList)

		if err != nil {
			// @todo do something useful here
		}

		trackedUsers := make([]adapters.LocalUser, len(twitterUsers))

		for index, user := range twitterUsers {
			trackedUsers[index] = adapters.NewUser(user).ToLocalUser()
		}

		respondWithSuccess(writer, http.StatusOK, trackedUsers)
	}
}
