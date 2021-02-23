package handlers

import (
	"encoding/json"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/mux"
	"insights-api/adapters"
	"insights-api/mocks"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFollowedUsersHandler_Success(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost, "/users/tracking", nil)
	rr := httptest.NewRecorder()

	mockService := mocks.MockUser{
		MockGetFollowingList: func() ([]twitter.User, error) {
			mockTwitterUsers := mocks.CreateMockTwitterUsers(3)
			return mockTwitterUsers, nil
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/users/tracking", FollowedUsersHandler(mockService))
	router.ServeHTTP(rr, request)

	var result []adapters.LocalUser
	err := json.Unmarshal([]byte(rr.Body.String()), &result)

	log.Println(rr.Body.String())

	if err != nil {
		t.Fatalf("Unexpected error unmarshalling HTTP response: %v", err.Error())
	}

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code: %d but got: %d", http.StatusOK, rr.Code)
	}

	if len(result) != 3 {
		t.Errorf("Expected to be following: %d users but got: %d users", 3, len(result))
	}

	log.Println(result)
}
