package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/muser"
)

//isUser - check wheter user logged in and is in the user base
func isUser(r *http.Request) (bool, error) {
	mail, _ := muser.GetMailByToken(r)
	if mail == "" {
		return false, errors.New("Authorization required")
	}
	return cdb.Misc().UserExist(r.Context(), mail)
}

//GetCurrentUser - get currently loggined user from auth header
func GetCurrentUser(r *http.Request) (users.User, error) {
	mail, _ := muser.GetMailByToken(r)
	if mail == "" {
		return users.User{}, errors.New("Authorization required")
	}
	return cdb.Misc().UserByEmail(r.Context(), mail)
}

const AdminRoleName = "Admin"

func isAdmin(r *http.Request) (bool, error) {
	user, err := GetCurrentUser(r)
	if err != nil {
		return false, err
	}
	return user.Role == AdminRoleName, nil
}

//UserCheck - handler wraper to check access rights
func UserCheck(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isHe, err := isUser(r); err != nil || !isHe {
			http.Error(w, "User access required", http.StatusForbidden)
			return // don't call original handler
		}
		h.ServeHTTP(w, r)
	})
}

//AdminCheck - handler wraper to check access rights
func AdminCheck(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isHe, err := isAdmin(r); err != nil || !isHe {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return // don't call original handler
		}
		h.ServeHTTP(w, r)
	})
}

//UserCheckHandler - returns current user if it exist
func UserCheckHandler(w http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(r)
	if err != nil {
		http.Error(w, "Authorization error: "+err.Error(), http.StatusForbidden)
		return // don't call original handler
	}
	b, err := json.Marshal(user)
	if err != nil {
		log.Println("Can't json.Marshal(user) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(b))
}
