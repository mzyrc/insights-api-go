package handlers

import (
	"encoding/json"
	"insights-api/adapters"
	"net/http"
)

func FollowedUsersHandler(service TwitterUserService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		followedUsers, _ := service.GetFollowingList()

		users := make([]adapters.LocalUser, len(followedUsers))

		for index, user := range followedUsers {
			users[index] = adapters.NewUser(user).ToLocalUser()
		}

		response, _ := json.Marshal(users)

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(response)
	}
}
