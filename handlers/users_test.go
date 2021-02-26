package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/mux"
	"insights-api/adapters"
	"insights-api/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsersHandler_Success(t *testing.T) {
	mockTwitterUsers := mocks.CreateMockTwitterUsers(1)
	postBody := map[string]string{
		"screen_name": mockTwitterUsers[0].ScreenName,
	}

	jsonBody, _ := json.Marshal(postBody)

	request, _ := http.NewRequest(http.MethodPost, "/users/query", bytes.NewBuffer(jsonBody))
	rr := httptest.NewRecorder()

	mockService := mocks.MockUser{
		MockSearch: func(screenName string) (users []twitter.User, err error) {
			return mockTwitterUsers, nil
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/users/query", UserSearchHandler(mockService))
	router.ServeHTTP(rr, request)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code: %d but got: %d", http.StatusOK, rr.Code)
	}

	var data []adapters.LocalUser
	err := json.Unmarshal([]byte(rr.Body.String()), &data)

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling response: %v", err.Error())
	}

	expected := adapters.LocalUser{
		ID:                       mockTwitterUsers[0].IDStr,
		Name:                     mockTwitterUsers[0].Name,
		ScreenName:               mockTwitterUsers[0].ScreenName,
		FollowersCount:           mockTwitterUsers[0].FollowersCount,
		Verified:                 mockTwitterUsers[0].Verified,
		SignedUpAt:               mockTwitterUsers[0].CreatedAt,
		TweetCount:               mockTwitterUsers[0].StatusesCount,
		ProfileImageThumbnailURL: mockTwitterUsers[0].ProfileImageURL,
		ProfileImageFullSizeURL:  mockTwitterUsers[0].ProfileImageURL,
	}

	if data[0] != expected {
		t.Errorf("Expected %v but got %v", expected, data[0])
	}
}