package handlers

import (
	"encoding/json"
	"github.com/dghubble/go-twitter/twitter"
	"insights-api/adapters"
	"net/http"
)

type TwitterUserService interface {
	Search(screenName string) (users []twitter.User, err error)
	GetFollowingList() ([]twitter.User, error)
	GetUser(userId int64) (*twitter.User, error)
	GetUsers(userIdList []int64) ([]twitter.User, error)
}

type userSearchPostBody struct {
	ScreenName string `json:"screen_name"`
}

func UserSearchHandler(service TwitterUserService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		decoder := json.NewDecoder(request.Body)
		var searchPostBody userSearchPostBody

		postBodyErr := decoder.Decode(&searchPostBody)

		if postBodyErr != nil {
			respondWithError(writer, http.StatusBadRequest, "Invalid post body")
			return
		}

		twitterUsers, err := service.Search(searchPostBody.ScreenName)

		if err != nil {
			respondWithError(writer, http.StatusBadRequest, err.Error())
			return
		}

		users := make([]adapters.LocalUser, len(twitterUsers))

		for index, user := range twitterUsers {
			userAdapter := adapters.NewUser(user)
			users[index] = userAdapter.ToLocalUser()
		}

		respondWithSuccess(writer, http.StatusOK, users)
	}
}
