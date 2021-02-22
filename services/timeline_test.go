package services

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

	service := Timeline{client: mockTimeLineService}

	_, err := service.GetTimeLine(1)

	if err != nil {
		t.Fatal(err)
	}
}

func TestTimelineService_GetTimeLine_Error(t *testing.T) {
	testCases := []struct {
		description    string
		mockTimelineFn func(config TimelineConfig) ([]twitter.Tweet, *http.Response, error)
		expectedError  string
	}{
		{
			description: "Should return an error when the user does not exist",
			mockTimelineFn: func(config TimelineConfig) ([]twitter.Tweet, *http.Response, error) {
				mockHttpResponse := http.Response{StatusCode: 404}
				return nil, &mockHttpResponse, errors.New("mock error message")
			},
			expectedError: "user does not exist",
		},
		{
			description: "Should return an unexpected error occurred error when the error is a not found error",
			mockTimelineFn: func(config TimelineConfig) ([]twitter.Tweet, *http.Response, error) {
				mockHttpResponse := http.Response{StatusCode: 500}
				return nil, &mockHttpResponse, errors.New("mock error message")
			},
			expectedError: "unexpected error occurred requesting timeline: mock error message",
		},
		{
			description: "Should return an unexpected error occurred error when an HTTP response is not present",
			mockTimelineFn: func(config TimelineConfig) ([]twitter.Tweet, *http.Response, error) {
				return nil, nil, errors.New("mock error message")
			},
			expectedError: "unexpected error occurred requesting timeline: mock error message",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			mockTimeLineService := MockTimelineProcessor{
				mockGetTimeLine: testCase.mockTimelineFn,
			}

			service := Timeline{client: mockTimeLineService}

			_, err := service.GetTimeLine(1)

			if err.Error() != testCase.expectedError {
				t.Fatalf("Expected %v but got %v", testCase.expectedError, err)
			}
		})
	}
}
