package orders_test

import (
	"context"
	"testing"

	"github.com/hromov/jevelina/domain/orders"
	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/mocks"
	"github.com/stretchr/testify/mock"
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
		{
			Email:        "user5@email.com",
			Distribution: 0.0,
		},
	}
	cs := &mocks.ContactsService{}
	ls := &mocks.LeadsService{}
	us := &mocks.UsersService{}
	us.On("List", mock.Anything).Return(users, nil)
	s := orders.NewService(cs, ls, us)

	resMap := make(map[string]int)
	for i := 0; i < 200; i++ {
		user, err := s.GetRandomUser(context.TODO())
		require.NoError(t, err)
		require.NotNil(t, user)
		resMap[user.Email] += 1
	}
	t.Log(resMap)
}
