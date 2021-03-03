package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"insights-api/dao"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type MockUntrackUserDAO struct {
	MockRemove func(userId int, twitterUserId int64) error
}

func (m MockUntrackUserDAO) Remove(userId int, twitterUserId int64) error {
	return m.MockRemove(userId, twitterUserId)
}

func TestUntrackUserHandler(t *testing.T) {
	testCases := []struct {
		description          string
		userId               string
		mockRemove           func(userId int, twitterUserId int64) error
		expectedStatusCode   int
		expectedResponseBody httpResponsePayload
	}{
		{
			description: "Should respond with status 200 OK when successfully removing a track on a user",
			userId:      "1",
			mockRemove: func(userId int, twitterUserId int64) error {
				return nil
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: httpResponsePayload{Data: ""},
		},
		{
			description: "Should respond with status 400 Bad Request when an invalid user id is supplied in the URL",
			userId:      "an invalid user id",
			mockRemove: func(userId int, twitterUserId int64) error {
				return nil
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: httpResponsePayload{Error: "Invalid user id in URL"},
		},
		{
			description: "Should respond with 404 Not Found when the supplied user id is not tracked",
			userId:      "100",
			mockRemove: func(userId int, twitterUserId int64) error {
				return dao.UserNotExistsError{Id: twitterUserId}
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: httpResponsePayload{Error: "User id 100 is not tracked"},
		},
		{
			description: "Should respond with 500 Internal Server error when an unknown error occurs",
			userId:      "100",
			mockRemove: func(userId int, twitterUserId int64) error {
				return errors.New("some unknown error")
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: httpResponsePayload{Error: "Unknown error occurred"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			untrackUserDAO := MockUntrackUserDAO{MockRemove: testCase.mockRemove}

			requestURL := fmt.Sprintf("/user/%s/untrack", testCase.userId)
			request, _ := http.NewRequest(http.MethodDelete, requestURL, nil)
			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/user/{userId}/untrack", UntrackUserHandler(untrackUserDAO))
			router.ServeHTTP(rr, request)

			if rr.Code != testCase.expectedStatusCode {
				t.Errorf("Expected status code: %d but got status code %d", testCase.expectedStatusCode, rr.Code)
			}

			responseBody := rr.Body.String()

			var actual httpResponsePayload
			json.Unmarshal([]byte(responseBody), &actual)

			if !reflect.DeepEqual(actual, testCase.expectedResponseBody) {
				t.Errorf("Expected response: %q but got: %q", testCase.expectedResponseBody, responseBody)
			}
		})
	}
}
