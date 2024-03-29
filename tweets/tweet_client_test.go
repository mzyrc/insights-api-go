package tweets

import (
	"errors"
	"github.com/dghubble/go-twitter/twitter"
	"net/http"
	"testing"
)

type MockTimelineProcessor struct {
	mockGetTimeLine func(config TimelineConfig) ([]twitter.Tweet, *http.Response, error)
}

func (m MockTimelineProcessor) UserTimeline(config TimelineConfig) ([]twitter.Tweet, *http.Response, error) {
	return m.mockGetTimeLine(config)
}

func TestTimelineService_GetTimeLine(t *testing.T) {
	mockTimeLineService := MockTimelineProcessor{
		mockGetTimeLine: func(config TimelineConfig) ([]twitter.Tweet, *http.Response, error) {
			mockHttpResponse := http.Response{StatusCode: 200}
			mockTweet := twitter.Tweet{
				ID:   123456,
				Text: "mock tweet text",
			}

			tweets := make([]twitter.Tweet, 1)
			tweets[0] = mockTweet

			return tweets, &mockHttpResponse, nil
		},
	}

	service := tweetClient{client: mockTimeLineService}

	_, err := service.GetTimeLine(TimelineConfig{UserId: 1})

	if err != nil {
		t.Fatal(err)
	}
}

func TestTimelineService_GetTimeLine_Error(t *testing.T) {
	testCases := []struct {
		description       string
		mockTimelineFn    func(config TimelineConfig) ([]twitter.Tweet, *http.Response, error)
		tweetClientConfig TimelineConfig
		expectedError     string
	}{
		{
			description: "Should return an error when a user id property is not specified on the config",
			mockTimelineFn: func(config TimelineConfig) ([]twitter.Tweet, *http.Response, error) {
				return nil, nil, nil
			},
			tweetClientConfig: TimelineConfig{UserId: 0},
			expectedError:     "must specify a a UserID property in the config",
		},
		{
			description: "Should return an error when the user does not exist",
			mockTimelineFn: func(config TimelineConfig) ([]twitter.Tweet, *http.Response, error) {
				mockHttpResponse := http.Response{StatusCode: 404}
				return nil, &mockHttpResponse, errors.New("mock error message")
			},
			tweetClientConfig: TimelineConfig{UserId: 1},
			expectedError:     "user does not exist",
		},
		{
			description: "Should return an unexpected error occurred error when the error is a not found error",
			mockTimelineFn: func(config TimelineConfig) ([]twitter.Tweet, *http.Response, error) {
				mockHttpResponse := http.Response{StatusCode: 500}
				return nil, &mockHttpResponse, errors.New("mock error message")
			},
			tweetClientConfig: TimelineConfig{UserId: 1},
			expectedError:     "unexpected error occurred requesting TweetService: mock error message",
		},
		{
			description: "Should return an unexpected error occurred error when an HTTP response is not present",
			mockTimelineFn: func(config TimelineConfig) ([]twitter.Tweet, *http.Response, error) {
				return nil, nil, errors.New("mock error message")
			},
			tweetClientConfig: TimelineConfig{UserId: 1},
			expectedError:     "unexpected error occurred requesting TweetService: mock error message",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mockTimeLineService := MockTimelineProcessor{
				mockGetTimeLine: testCase.mockTimelineFn,
			}

			service := tweetClient{client: mockTimeLineService}

			_, err := service.GetTimeLine(testCase.tweetClientConfig)

			if err.Error() != testCase.expectedError {
				t.Fatalf("Expected %v but got %v", testCase.expectedError, err)
			}
		})
	}
}
