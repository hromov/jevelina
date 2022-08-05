package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hromov/jevelina/api"
	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type roleCreateTest struct {
	name    string
	req     string
	resp    users.Role
	respErr error
	code    int
}

var createRoleTests = []roleCreateTest{
	{
		name: "succes role creation",
		req:  `{"Priority": 5, "Role": "some name"}`,
		resp: users.Role{
			ID:        uint8(1),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: time.Time{},
			Priority:  5,
			Role:      "some name",
		},
		respErr: nil,
		code:    http.StatusAccepted,
	},
	{
		name:    "json failed role creation",
		req:     `{"Priority": "5", "Role": "some name"}`,
		resp:    users.Role{},
		respErr: nil,
		code:    http.StatusBadRequest,
	},
	{
		name:    "db failed role creation",
		req:     `{"Priority": 5, "Role": "some name"}`,
		resp:    users.Role{},
		respErr: errors.New("some error"),
		code:    http.StatusInternalServerError,
	},
}

func TestCreateRole(t *testing.T) {
	for _, tc := range createRoleTests {
		req, err := http.NewRequest(
			"POST",
			"/roles",
			bytes.NewBuffer([]byte(
				[]byte(tc.req),
			)),
		)
		require.NoError(t, err)

		us := &mocks.UsersService{}
		us.On("CreateRole", mock.Anything, mock.Anything).Return(tc.resp, tc.respErr)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(api.CreateRole(us))

		handler.ServeHTTP(rr, req)

		require.Equal(t, tc.code, rr.Code)

		//we test result body only for success cases
		if tc.code < 400 {
			b, _ := json.Marshal(tc.resp)
			require.Equal(t, string(b), rr.Body.String())
		}
	}
}
