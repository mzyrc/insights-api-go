package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"insights-api/dao"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockTrackedUserDAO struct {
	mockCreate func(userId int, twitterUserId int64) error
}

func (m MockTrackedUserDAO) Create(userId int, twitterUserId int64) error {
	return m.mockCreate(userId, twitterUserId)
}

func TestTrackUserHandler(t *testing.T) {
	testCases := []struct {
		description          string
		postBody             interface{}
		mockCreate           func(userId int, twitterUserId int64) error
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			description: "Should respond with 200 OK when creating a new tracking succeeds",
			postBody:    map[string]string{"user_id": "1234567890"},
			mockCreate: func(userId int, twitterUserId int64) error {
				return nil
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "",
		},
		{
			description: "Should respond with 409 Conflict when a tracking for the specified user exists",
			postBody:    map[string]string{"user_id": "1234567890"},
			mockCreate: func(userId int, twitterUserId int64) error {
				return dao.UserExistsError{Id: twitterUserId}
			},
			expectedStatusCode:   http.StatusConflict,
			expectedResponseBody: "tracked user: 1234567890 already tracked",
		},
		{
			description: "Should respond with a 400 Bad Request if the post body is invalid",
			postBody:    "invalid post body",
			mockCreate: func(userId int, twitterUserId int64) error {
				return nil
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "Invalid post body received",
		},
		{
			description: "Should respond with a 400 Bad Request when an invalid user_id is specified in the post body",
			postBody:    map[string]string{"user_id": "invalid user id"},
			mockCreate: func(userId int, twitterUserId int64) error {
				return nil
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "Invalid user_id in post body",
		},
		{
			description: "Should respond with 500 if an unknown error has occurred",
			postBody:    map[string]string{"user_id": "1234567890"},
			mockCreate: func(userId int, twitterUserId int64) error {
				return errors.New("some unknown error")
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: "An unknown error occurred",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mockTrackedUserDAO := MockTrackedUserDAO{mockCreate: testCase.mockCreate}

			jsonBody, _ := json.Marshal(testCase.postBody)
			request, _ := http.NewRequest(http.MethodPost, "/user/track/new", bytes.NewBuffer(jsonBody))
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/user/track/new", TrackUserHandler(mockTrackedUserDAO))
			router.ServeHTTP(rr, request)

			if rr.Code != testCase.expectedStatusCode {
				t.Errorf("Expected status code: %d but got: %d", testCase.expectedStatusCode, rr.Code)
			}

			responseBody := rr.Body.String()

			if responseBody != testCase.expectedResponseBody {
				t.Errorf("Expected body: %q but got: %q", testCase.expectedResponseBody, responseBody)
			}
		})
	}
}
