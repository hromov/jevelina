package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/models"
	"github.com/hromov/jevelina/cdb/orders"
)

const randomUserEmail = "random@random.org"

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/orders" {
		http.NotFound(w, r)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Only POST is allowed", http.StatusForbidden)
		return
	}

	c := new(models.CreateLeadReq)
	if err := json.NewDecoder(r.Body).Decode(c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if c.UserEmail == "" || c.UserHash == "" {
		http.Error(w, "User Email and Hash are required", http.StatusBadRequest)
		return
	}
	user, err := cdb.Misc().UserByEmail(c.UserEmail)
	if err != nil || user == nil {
		http.Error(w, "Cant find user with email: "+c.UserEmail, http.StatusBadRequest)
		return
	}

	if user.Hash != c.UserHash {
		http.Error(w, "Wrong user-hash values", http.StatusForbidden)
		return
	}
	if user.Email == randomUserEmail {
		user, err = getRandomUser()
		if err != nil {
			log.Println("Can't get random user error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
	}
	contact, err := orders.CreateOrGetContact(c, user)
	if err != nil {
		log.Println("Can't create contact error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	lead, err := orders.CreateLead(c, contact)
	if err != nil {
		log.Println("Can't create lead error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	if err = orders.CreateTask(c, lead); err != nil {
		log.Println("Can't create task error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(lead)
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

func getRandomUser() (*models.User, error) {
	users, err := cdb.Misc().Users()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("No users found in base")
	}

	rand.Seed(time.Now().UnixNano())
	var user *models.User
	appropriateSeen := false
	for user == nil {
		for _, u := range users {
			if u.Distribution == 0.0 {
				continue
			}
			appropriateSeen = true
			r := rand.Float32()
			// log.Printf("checking %+v, rand = %.2f, good = %v\n", u, r, u.Distribution >= r)
			if u.Distribution >= r {
				return &u, nil
			}
		}
		if !appropriateSeen {
			return nil, errors.New("No user with distribution more then 0 was found")
		}
	}
	return nil, errors.New("should never be called")
}
