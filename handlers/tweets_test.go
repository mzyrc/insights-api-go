package handlers

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockTimelineProcessor struct {
	mockGetTimeline func(userId int64) ([]twitter.Tweet, error)
}

func (m MockTimelineProcessor) GetTimeLine(userId int64) ([]twitter.Tweet, error) {
	return m.mockGetTimeline(userId)
}

func TestTweetsHandler(t *testing.T) {
	testCases := []struct {
		description        string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			description:        "Should return a list of tweets",
			expectedStatusCode: http.StatusOK,
			expectedBody:       "[{\"coordinates\":null,\"created_at\":\"\",\"current_user_retweet\":null,\"entities\":null,\"favorite_count\":0,\"favorited\":false,\"filter_level\":\"\",\"id\":0,\"id_str\":\"\",\"in_reply_to_screen_name\":\"\",\"in_reply_to_status_id\":0,\"in_reply_to_status_id_str\":\"\",\"in_reply_to_user_id\":0,\"in_reply_to_user_id_str\":\"\",\"lang\":\"\",\"possibly_sensitive\":false,\"quote_count\":0,\"reply_count\":0,\"retweet_count\":0,\"retweeted\":false,\"retweeted_status\":null,\"source\":\"\",\"scopes\":null,\"text\":\"\",\"full_text\":\"\",\"display_text_range\":[0,0],\"place\":null,\"truncated\":false,\"user\":null,\"withheld_copyright\":false,\"withheld_in_countries\":null,\"withheld_scope\":\"\",\"extended_entities\":null,\"extended_tweet\":null,\"quoted_status_id\":0,\"quoted_status_id_str\":\"\",\"quoted_status\":null}]\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/user/12345/tweets", nil)

			rr := httptest.NewRecorder()

			mockService := MockTimelineProcessor{
				mockGetTimeline: func(userId int64) ([]twitter.Tweet, error) {
					mockTweets := make([]twitter.Tweet, 1)
					mockTweet := twitter.Tweet{}

					mockTweets[0] = mockTweet
					return mockTweets, nil
				},
			}

			router := mux.NewRouter()
			router.HandleFunc("/user/{userId}/tweets", TweetsHandler(mockService))
			router.ServeHTTP(rr, request)

			statusCode := rr.Code
			body := rr.Body.String()

			if statusCode != testCase.expectedStatusCode {
				t.Errorf("Expected %d but got %d", testCase.expectedStatusCode, statusCode)
			}

			if result := strings.Compare(body, testCase.expectedBody); result != 0 {
				t.Errorf("Expected %q but got %q", testCase.expectedBody, body)
			}
		})
	}
}
