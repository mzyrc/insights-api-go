package services

import (
	"github.com/dghubble/go-twitter/twitter"
	"insights-api/mocks"
	"net/http"
	"testing"
)

func TestUser_GetFollowingList_SuccessNoNextPage(t *testing.T) {
	mockFollowedUsers := mocks.CreateMockTwitterUsers(3)

	mockUserProcessor := mocks.MockTwitterUserClient{
		MockGetFriends: func(nextPageId int64) (*twitter.Friends, *http.Response, error) {
			mockFriends := twitter.Friends{
				Users:      mockFollowedUsers,
				NextCursor: 0,
			}

			return &mockFriends, nil, nil
		},
	}

	service := NewUser(&mockUserProcessor)

	result, err := service.GetFollowingList()

	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}

	actualNumberOfFollowedUsers := len(result)
	expectedNumberOfFollowedUsers := len(mockFollowedUsers)

	if actualNumberOfFollowedUsers != expectedNumberOfFollowedUsers {
		t.Errorf("Expected to be following %d users but got %d users", expectedNumberOfFollowedUsers, actualNumberOfFollowedUsers)
	}
}

func TestUser_GetFollowingList_SuccessWithNextPage(t *testing.T) {
	mockUserProcessor := mocks.MockTwitterUserClient{
		MockGetFriends: func(nextPageId int64) (*twitter.Friends, *http.Response, error) {
			var mockFollowedUsers []twitter.User

			switch nextPageId {
			case 0:
				mockFollowedUsers = mocks.CreateMockTwitterUsers(10)
				nextPageId = 1
				break
			case 1:
				mockFollowedUsers = mocks.CreateMockTwitterUsers(10)
				nextPageId = 2
			case 2:
				mockFollowedUsers = mocks.CreateMockTwitterUsers(5)
				nextPageId = 0
			}

			mockFriendsResponse := twitter.Friends{
				Users:      mockFollowedUsers,
				NextCursor: nextPageId,
			}

			mockHttpResponse := http.Response{StatusCode: http.StatusOK}

			return &mockFriendsResponse, &mockHttpResponse, nil
		},
	}

	service := NewUser(&mockUserProcessor)

	result, err := service.GetFollowingList()

	if err != nil {
		t.Errorf("Unexpected error: %v", err.Error())
	}

	expectedNumberOfFollowedUsers := 25

	if len(result) != expectedNumberOfFollowedUsers {
		t.Errorf("Expected to be following %d users but actually following %d users", expectedNumberOfFollowedUsers, len(result))
	}
}
