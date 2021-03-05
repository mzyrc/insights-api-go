package mocks

type MockTweetsSync struct {
	MockStoreTweetsByUser func(userId int64) error
}

func (m MockTweetsSync) StoreTweetsByUser(userId int64) error {
	return m.MockStoreTweetsByUser(userId)
}
