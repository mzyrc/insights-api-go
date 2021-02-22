package adapters

import (
	"github.com/dghubble/go-twitter/twitter"
	"strings"
)

type LocalUser struct {
	ID                       string `json:"id"`
	Name                     string `json:"name"`
	ScreenName               string `json:"screen_name"`
	ProfileImageThumbnailURL string `json:"profile_image_thumbnail_url"`
	ProfileImageFullSizeURL  string `json:"profile_image_full_size_url"`
	TweetCount               int    `json:"tweet_counts"`
	FollowersCount           int    `json:"followers_count"`
	Verified                 bool   `json:"verified"`
	SignedUpAt               string `json:"signed_up_at"`
}

type User struct {
	user twitter.User
}

func NewUser(user twitter.User) *User {
	return &User{user: user}
}

func (u User) ToLocalUser() LocalUser {
	return LocalUser{
		ID: u.user.IDStr,
		Name:                     u.user.Name,
		ScreenName:               u.user.ScreenName,
		ProfileImageThumbnailURL: u.user.ProfileImageURL,
		ProfileImageFullSizeURL:  strings.Replace(u.user.ProfileImageURL, "_normal", "_400x400", 1),
		TweetCount:               u.user.StatusesCount,
		FollowersCount:           u.user.FollowersCount,
		Verified:                 u.user.Verified,
		SignedUpAt:               u.user.CreatedAt,
	}
}
