package handlers

import (
	"insights-api/user"
	"net/http"
)

type trackedUsersDAO interface {
	GetUsers(userId int) ([]int64, error)
}

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

		twitterUsers, err := service.GetUsers(userIdList)

		if err != nil {
			// @todo do something useful here
		}

		trackedUsers := make([]user.LocalUser, len(twitterUsers))

		for index, twitterUser := range twitterUsers {
			trackedUsers[index] = user.NewUserAdapter(twitterUser).ToLocalUser()
		}

		respondWithSuccess(writer, http.StatusOK, trackedUsers)
	}
}
