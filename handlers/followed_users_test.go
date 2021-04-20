package handlers

import (
	"errors"
	"github.com/dghubble/go-twitter/twitter"
	"insights-api/mocks"
	"insights-api/user"
	"log"
	"reflect"
	"testing"
)

type mockUserLookUpService struct {
	MockGetUsers func(userIdList []int64) ([]twitter.User, error)
}

func (m mockUserLookUpService) GetUsers(userIdList []int64) ([]twitter.User, error) {
	return m.MockGetUsers(userIdList)
}

func TestGetTrackedUsersHandler(t *testing.T) {

}

func TestGetTrackedUsersSuccess(t *testing.T) {
	mockTwitterUsers := mocks.CreateMockTwitterUsers(1)
	mockService := mockUserLookUpService{
		MockGetUsers: func(userIdList []int64) ([]twitter.User, error) {
			return mockTwitterUsers, nil
		},
	}
	trackedUsers, err := getTrackedUsers(mockService, []int64{})

	if err != nil {
		t.Fatalf("Test failed with %v", err)
	}

	expected := user.NewUserAdapter(mockTwitterUsers[0]).ToLocalUser()

	if !reflect.DeepEqual(expected, trackedUsers[0]) {
		t.Errorf("Expected: %v but got got: %v", mockTwitterUsers, trackedUsers)
	}
}

func TestGetTrackedUsersReturnsAndErrorWhenTheLookUpFails(t *testing.T) {
	mockService := mockUserLookUpService{
		MockGetUsers: func(userIdList []int64) ([]twitter.User, error) {
			return nil, errors.New("Mock error")
		},
	}
	_, err := getTrackedUsers(mockService, []int64{})

	if err == nil {
		t.Fatalf("Test failed with %v", err)
	}

	log.Println(err)
}
