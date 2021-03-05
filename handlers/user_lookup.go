package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"insights-api/user"
	"net/http"
	"strconv"
)

func UserLookupHandler(service TwitterUserService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		requestVariables := mux.Vars(request)
		userId := requestVariables["userId"]

		twitterUserId, err := strconv.ParseInt(userId, 10, 64)

		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			writer.Write([]byte("Invalid user id in URL"))
			return
		}

		twitterUser, serviceErr := service.GetUser(twitterUserId)

		if serviceErr != nil {
			writer.WriteHeader(http.StatusServiceUnavailable)
			writer.Write([]byte("Unexpected error requesting the user"))
			return
		}

		response, _ := json.Marshal(user.NewUserAdapter(*twitterUser).ToLocalUser())

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(response)
	}
}
