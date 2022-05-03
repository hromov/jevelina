package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/harlow/authtoken"
	"github.com/hromov/jevelina/base"
	"github.com/hromov/jevelina/cdb/models"
)

const accesPath = "https://www.googleapis.com/oauth2/v3/tokeninfo?access_token="

type googleAuthResponse struct {
	Aud           string
	ExpiresIn     string `json:"expires_in"`
	Scope         string
	Email         string
	EmailVerified string `json:"email_verified"`
}

// GetMailByToken -
func GetMailByToken(r *http.Request) (string, error) {
	token, err := authtoken.FromRequest(r)
	if err != nil {
		log.Printf("token from request error: %v", err)
		return "", err
	}
	if token == "expired" {
		return "", errors.New("expired")
	}
	resp, err := http.Get(accesPath + token)
	if err != nil {
		log.Printf("token check error: %v", err)
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll error: %v", err)
		return "", err
	}

	answer := new(googleAuthResponse)
	err = json.Unmarshal(body, answer)
	if err != nil {
		log.Printf("json.Unmarashal error: %v", err)
		return "", err
	}

	expires, err := strconv.Atoi(answer.ExpiresIn)
	if err != nil {
		log.Printf("strconv.Atoi(expires) error: %v", err)
		return "", err
	}
	if expires > 0 {
		return answer.Email, nil
	}
	return "", errors.New("Token expired")
}

//isUser - check wheter user logged in and is in the user base
func isUser(r *http.Request) (bool, error) {
	mail, _ := GetMailByToken(r)
	if mail == "" {
		return false, errors.New("Authorization required")
	}
	return base.GetDB().Misc().UserExist(mail)
}

//GetCurrentUser - get currently loggined user from auth header
func GetCurrentUser(r *http.Request) (*models.User, error) {
	mail, _ := GetMailByToken(r)
	if mail == "" {
		return nil, errors.New("Authorization required")
	}
	user, err := base.GetDB().Misc().UserByEmail(mail)
	if err != nil {
		return nil, err
	}
	return user, nil
}

const AdminRoleName = "Admin"

func isAdmin(r *http.Request) (bool, error) {
	user, err := GetCurrentUser(r)
	if err != nil {
		return false, err
	}
	return user.Role.Role == AdminRoleName, nil
}

//UserCheck - handler wraper to check access rights
func UserCheck(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isHe, err := isUser(r); err != nil || !isHe {
			http.Error(w, "User access required", http.StatusForbidden)
			return // don't call original handler
		}
		h.ServeHTTP(w, r)
	})
}

//AdminCheck - handler wraper to check access rights
func AdminCheck(h http.HandlerFunc) http.HandlerFunc {
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
	fmt.Fprintf(w, string(b))
	// it said that its already ok now
	// w.WriteHeader(http.StatusOK)
	return
}
