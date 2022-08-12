package auth_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/http/rest/auth"
	"github.com/hromov/jevelina/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {}

func TestUserCheck(t *testing.T) {
	us := &mocks.UsersService{}
	user := users.User{ID: 1, Role: "User", Email: "some@mail.com"}
	us.On("GetByEmail", mock.Anything, mock.Anything).Return(user, nil)
	tv := &mocks.TokenService{}
	tv.On("GetMailByToken", mock.Anything).Return("some@mail.com", nil)
	as := auth.NewService(us, tv)
	router := mux.NewRouter()
	router.HandleFunc("/", dummyHandler).Methods("GET")
	router.Use(as.UserCheck)
	require.NotNil(t, router)

	req, err := http.NewRequest(
		"GET",
		"/",
		nil,
	)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}

func TestNotAUserCheck(t *testing.T) {
	us := &mocks.UsersService{}
	us.On("GetByEmail", mock.Anything, mock.Anything).Return(users.User{}, errors.New("Not a user"))
	tv := &mocks.TokenService{}
	tv.On("GetMailByToken", mock.Anything).Return("some@mail.com", nil)
	as := auth.NewService(us, tv)
	router := mux.NewRouter()
	router.HandleFunc("/", dummyHandler).Methods("GET")
	router.Use(as.UserCheck)
	require.NotNil(t, router)

	req, err := http.NewRequest(
		"GET",
		"/",
		nil,
	)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestNotTokenUserCheck(t *testing.T) {
	us := &mocks.UsersService{}
	tv := &mocks.TokenService{}
	tv.On("GetMailByToken", mock.Anything).Return("", errors.New("some wrong token"))
	as := auth.NewService(us, tv)
	router := mux.NewRouter()
	router.HandleFunc("/", dummyHandler).Methods("GET")
	router.Use(as.UserCheck)
	require.NotNil(t, router)

	req, err := http.NewRequest(
		"GET",
		"/",
		nil,
	)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestNotAdminCheck(t *testing.T) {
	us := &mocks.UsersService{}
	user := users.User{ID: 1, Role: "User", Email: "some@mail.com"}
	us.On("GetByEmail", mock.Anything, mock.Anything).Return(user, nil)
	tv := &mocks.TokenService{}
	tv.On("GetMailByToken", mock.Anything).Return("some@mail.com", nil)
	as := auth.NewService(us, tv)
	router := mux.NewRouter()
	router.HandleFunc("/", dummyHandler).Methods("GET")
	router.Use(as.UserCheck)
	router.Use(as.AdminCheck)
	require.NotNil(t, router)

	req, err := http.NewRequest(
		"GET",
		"/",
		nil,
	)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusForbidden, rr.Code)
}

func TestAdminCheck(t *testing.T) {
	us := &mocks.UsersService{}
	user := users.User{ID: 1, Role: "Admin", Email: "some@mail.com"}
	us.On("GetByEmail", mock.Anything, mock.Anything).Return(user, nil)
	tv := &mocks.TokenService{}
	tv.On("GetMailByToken", mock.Anything).Return("some@mail.com", nil)
	as := auth.NewService(us, tv)
	router := mux.NewRouter()
	router.HandleFunc("/", dummyHandler).Methods("GET")
	router.Use(as.UserCheck)
	router.Use(as.AdminCheck)
	require.NotNil(t, router)

	req, err := http.NewRequest(
		"GET",
		"/",
		nil,
	)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
}
