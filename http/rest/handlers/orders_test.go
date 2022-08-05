package handlers

import (
	"testing"

	"github.com/hromov/jevelina/domain/users"
	"github.com/stretchr/testify/require"
)

func TestGetRandomUser(t *testing.T) {
	users := []users.User{
		{
			Email:        "user1@email.com",
			Distribution: 0.25,
		},
		{
			Email:        "user2@email.com",
			Distribution: 0.25,
		},
		{
			Email:        "user3@email.com",
			Distribution: 0.25,
		},
		{
			Email:        "user4@email.com",
			Distribution: 0.25,
		},
	}

	resMap := make(map[string]int)
	for i := 0; i < 200; i++ {
		user, err := getRandomUser(users)
		require.NoError(t, err)
		require.NotNil(t, user)
		resMap[user.Email] += 1
	}
	t.Log(resMap)
}
