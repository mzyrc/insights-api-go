package dao

import (
	"database/sql"
	"fmt"
	"github.com/lib/pq"
)

type UserExistsError struct {
	Id int64
}

func (u UserExistsError) Error() string {
	return fmt.Sprintf("tracked user: %d already tracked", u.Id)
}

type TrackedUser struct {
	db *sql.DB
}

func NewTrackedUserDAO(db *sql.DB) *TrackedUser {
	return &TrackedUser{db: db}
}

func (tu TrackedUser) Create(userId int, twitterUserId int64) error {
	insertSQL := "INSERT INTO usr_following (user_id, twitter_user_id) VALUES ($1, $2)"
	_, err := tu.db.Exec(insertSQL, userId, twitterUserId)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return UserExistsError{Id: twitterUserId}
		}
	}

	return nil
}