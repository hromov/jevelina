package orders_test

import (
	"context"
	"testing"

	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/mocks"
	"github.com/hromov/jevelina/useCases/orders"
	"github.com/hromov/jevelina/useCases/tasks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateForUser(t *testing.T) {
	u := users.User{
		ID: 1,
	}
	cs := &mocks.ContactsService{}
	cs.On("CreateOrGet", mock.Anything, mock.Anything).Return(contacts.Contact{
		ID:          234,
		Responsible: u,
	}, nil)
	ls := &mocks.LeadsService{}
	ls.On("Create", mock.Anything, leads.LeadData{
		ID:            0,
		ContactID:     234,
		ResponsibleID: u.ID,
		CreatedID:     u.ID,
	}).Return(leads.Lead{ID: 123, Responsible: u, Created: u}, nil)
	us := &mocks.UsersService{}
	ts := &mocks.TasksService{}
	ts.On("Create", mock.Anything, mock.Anything).Return(tasks.Task{}, nil)
	os := orders.NewService(cs, ls, us, ts)
	err := os.CreateForUser(context.TODO(), orders.Order{}, u)
	require.NoError(t, err)
}
