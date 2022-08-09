package auth

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/hromov/jevelina/domain/users"
	tokenvalidator "github.com/hromov/jevelina/utils/tokenValidator"
)

type KeyUser struct{}

const AdminRoleName = "Admin"

type Service interface {
	UserCheck(h http.Handler) http.Handler
	AdminCheck(h http.Handler) http.Handler
	UserCheckHandler() func(w http.ResponseWriter, r *http.Request)
}
type service struct {
	us users.Service
	tv tokenvalidator.Service
}

func NewService(us users.Service, tv tokenvalidator.Service) *service {
	return &service{us, tv}
}

//UserCheck - handler wraper to check access rights
func (s *service) UserCheck(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := s.getCurrentUser(r)
		if err != nil {
			log.Println("User check error: ", err.Error())
			http.Error(w, "User access required", http.StatusForbidden)
			return
		}

		if user.ID == 0 {
			log.Println("Empty user error")
			http.Error(w, "User access required", http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), KeyUser{}, user)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

//AdminCheck - handler wraper to check access rights
func (s *service) AdminCheck(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userValue := r.Context().Value(KeyUser{})
		user, ok := userValue.(users.User)
		if !ok || user.Role != AdminRoleName {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	})
}

//UserCheckHandler - returns current user if it exist
func (s *service) UserCheckHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := s.getCurrentUser(r)
		if err != nil {
			http.Error(w, "Authorization error: "+err.Error(), http.StatusForbidden)
			return
		}
		_ = json.NewEncoder(w).Encode(user)
	}
}

func (s *service) getCurrentUser(r *http.Request) (users.User, error) {
	mail, _ := s.tv.GetMailByToken(r)
	if mail == "" {
		return users.User{}, errors.New("Authorization required")
	}
	return s.us.GetByEmail(r.Context(), mail)
}
