package mocks

type MockTweetsSync struct {
	MockStoreTweetsForFirstTime func(userId int64)
}

func (m MockTweetsSync) StoreTweetsForFirstTime(userId int64) {
	m.MockStoreTweetsForFirstTime(userId)
}
