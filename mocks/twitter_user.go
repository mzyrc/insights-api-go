package mocks

import (
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	"strconv"
	"time"
)

func CreateMockTwitterUsers(numberToCreate int) []twitter.User {
	users := make([]twitter.User, numberToCreate)

	for index := 0; index < numberToCreate; index++ {
		userIDStr := fmt.Sprintf("1%d", index)
		userID, _ := strconv.ParseInt(userIDStr, 10, 64)
		createdAt := time.Date(2021, 1, index, 0, 0, 0, 0, &time.Location{})

		users[index] = twitter.User{
			CreatedAt:       createdAt.String(),
			FollowersCount:  10,
			ID:              userID,
			IDStr:           userIDStr,
			Name:            fmt.Sprintf("dummy name %d", numberToCreate),
			ProfileImageURL: "dummy image url",
			ScreenName:      fmt.Sprintf("dummy screen name %d", numberToCreate),
			StatusesCount:   0,
			Verified:        false,
		}
	}

	return users
}
