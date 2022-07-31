package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"crypto/rand"
	"math/big"
	"net/http"

	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/models"
	"github.com/hromov/jevelina/cdb/orders"
	d_users "github.com/hromov/jevelina/domain/users"
)

const randomUserEmail = "random@random.org"

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	user, err := cdb.Misc().UserByEmail(ctx, c.UserEmail)
	if err != nil || user.ID == 0 {
		http.Error(w, "Cant find user with email: "+c.UserEmail, http.StatusBadRequest)
		return
	}

	if user.Hash != c.UserHash {
		http.Error(w, "Wrong user-hash values", http.StatusForbidden)
		return
	}
	if user.Email == randomUserEmail {
		users, err := cdb.Misc().Users(ctx)
		if err != nil {
			http.Error(w, "Can't get users", http.StatusInternalServerError)
			return
		}
		if len(users) == 0 {
			http.Error(w, "Users length = 0", http.StatusInternalServerError)
			return
		}
		user, err = getRandomUser(users)
		if err != nil {
			log.Println("Can't get random user error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		} else {
			log.Println("randomUser selected to: ", user.Name)
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
	fmt.Fprint(w, string(b))
}

func getRandomUser(users []d_users.User) (d_users.User, error) {
	filtered := []d_users.User{}
	for _, u := range users {
		if u.Distribution > 0.0 {
			filtered = append(filtered, u)
		}
	}

	if len(filtered) == 0 {
		return d_users.User{}, errors.New("no good users found")
	}

	r, err := rand.Int(rand.Reader, big.NewInt(int64(len(filtered))))
	if err != nil {
		return d_users.User{}, err
	}
	return filtered[r.Int64()], nil
}
