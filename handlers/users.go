package handlers

import (
	"encoding/json"
	"github.com/dghubble/go-twitter/twitter"
	"insights-api/adapters"
	"net/http"
)

type UserService interface {
	Search(screenName string) (users []twitter.User, err error)
	GetFollowingList() ([]twitter.User, error)
}

type UserSearchParams struct {
	ScreenName string `json:"screen_name"`
}

func UsersHandler(service UserService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var postBody UserSearchParams

		postBodyErr := decoder.Decode(&postBody)

		if postBodyErr != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Invalid post body"))
			return
		}

		twitterUsers, err := service.Search(postBody.ScreenName)

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte(err.Error()))
			return
		}

		users := make([]adapters.LocalUser, len(twitterUsers))

		for index, user := range twitterUsers {
			userAdapter := adapters.NewUser(user)
			users[index] = userAdapter.ToLocalUser()
		}

		response, _ := json.Marshal(users)
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}

func FollowedUsersHandler(service UserService) http.HandlerFunc {
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
