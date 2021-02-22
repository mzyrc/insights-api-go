package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/mux"
	"insights-api/adapters"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockUserProcessor struct {
	mockSearch func(screenName string) (users []twitter.User, err error)
}

func (m MockUserProcessor) Search(screenName string) (users []twitter.User, err error) {
	return m.mockSearch(screenName)
}

func TestUsersHandler_Success(t *testing.T) {
	mockTwitterUser := twitter.User{
		ID:             123,
		IDStr:          "123",
		Name:           "mockName",
		ScreenName:     "mockScreenName",
		FollowersCount: 100,
		Verified:       false,
		CreatedAt:      "2021-01-01",
		StatusesCount:  100,
	}
	postBody := map[string]string{
		"screen_name": mockTwitterUser.ScreenName,
	}

	jsonBody, _ := json.Marshal(postBody)

	request, _ := http.NewRequest(http.MethodPost, "/users/query", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	mockService := MockUserProcessor{
		mockSearch: func(screenName string) (users []twitter.User, err error) {
			mockTwitterUsers := make([]twitter.User, 1)
			mockTwitterUsers[0] = mockTwitterUser

			return mockTwitterUsers, nil
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/users/query", UsersHandler(mockService))
	router.ServeHTTP(rr, request)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected %d but got %d", http.StatusOK, rr.Code)
	}

	var data []adapters.LocalUser
	err := json.Unmarshal([]byte(rr.Body.String()), &data)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling response: %v", err.Error())
	}

	expected := adapters.LocalUser{
		ID:                       mockTwitterUser.IDStr,
		Name:                     mockTwitterUser.Name,
		ScreenName:               mockTwitterUser.ScreenName,
		FollowersCount:           mockTwitterUser.FollowersCount,
		Verified:                 mockTwitterUser.Verified,
		SignedUpAt:               mockTwitterUser.CreatedAt,
		TweetCount:               mockTwitterUser.StatusesCount,
		ProfileImageThumbnailURL: "",
		ProfileImageFullSizeURL:  "",
	}

	if data[0] != expected {
		t.Errorf("Expected %v but got %v", expected, data[0])
	}
}
