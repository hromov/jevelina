package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"crypto/rand"
	"math/big"
	"net/http"

	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/misc"
	d_users "github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/storage/mysql"
	"github.com/hromov/jevelina/storage/mysql/dao/models"
	"github.com/hromov/jevelina/storage/mysql/dao/orders"
)

const randomUserEmail = "random@random.org"

type CreateLeadReq struct {
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description,omitempty"`

	ClientName  string `json:"clientname"`
	ClientEmail string `json:"clientemail,omitempty"`
	ClientPhone string `json:"clientphone,omitempty"`

	Source       string `json:"source,omitempty"`
	Product      string `json:"product,omitempty"`
	Manufacturer string `json:"manufacturer,omitempty"`

	UserEmail string `json:"user_email,omitempty"`
	UserHash  string `json:"user_hash,omitempty"`

	CID string
	UID string
	TID string

	UtmID       string
	UtmSource   string
	UtmMedium   string
	UtmCampaign string

	Domain string
}

func (c *CreateLeadReq) ToContactRequest(userID uint64) contacts.ContactRequest {
	cr := contacts.ContactRequest{
		Name:          c.ClientName,
		Email:         c.ClientEmail,
		Phone:         c.ClientPhone,
		ResponsibleID: userID,
		Analytics: misc.Analytics{
			CID:         c.CID,
			TID:         c.TID,
			UtmID:       c.UtmID,
			UtmSource:   c.UtmSource,
			UtmMedium:   c.UtmMedium,
			UtmCampaign: c.UtmCampaign,
			Domain:      c.Domain,
		},
	}
	if c.UID == "" {
		cr.Analytics.UID = c.ClientPhone
	} else {
		cr.Analytics.UID = c.UID
	}

	return cr
}

func Order(cs contacts.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.URL.Path != "/orders" {
			http.NotFound(w, r)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Only POST is allowed", http.StatusForbidden)
			return
		}

		c := CreateLeadReq{}
		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if c.UserEmail == "" || c.UserHash == "" {
			http.Error(w, "User Email and Hash are required", http.StatusBadRequest)
			return
		}
		user, err := mysql.Misc().UserByEmail(ctx, c.UserEmail)
		if err != nil || user.ID == 0 {
			http.Error(w, "Cant find user with email: "+c.UserEmail, http.StatusBadRequest)
			return
		}

		if user.Hash != c.UserHash {
			http.Error(w, "Wrong user-hash values", http.StatusForbidden)
			return
		}
		if user.Email == randomUserEmail {
			users, err := mysql.Misc().Users(ctx)
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

		contact, err := cs.CreateOrGet(r.Context(), c.ToContactRequest(user.ID))
		if err != nil {
			log.Println("Can't create contact error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		lead, err := orders.CreateLead(models.CreateLeadReq(c), contact)
		if err != nil {
			log.Println("Can't create lead error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		if err = orders.CreateTask(models.CreateLeadReq(c), lead); err != nil {
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
}

func Orders(cs contacts.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if r.URL.Path != "/orders" {
			http.NotFound(w, r)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Only POST is allowed", http.StatusForbidden)
			return
		}

		c := CreateLeadReq{}
		if err := json.NewDecoder(r.Body).Decode(c); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if c.UserEmail == "" || c.UserHash == "" {
			http.Error(w, "User Email and Hash are required", http.StatusBadRequest)
			return
		}
		user, err := mysql.Misc().UserByEmail(ctx, c.UserEmail)
		if err != nil || user.ID == 0 {
			http.Error(w, "Cant find user with email: "+c.UserEmail, http.StatusBadRequest)
			return
		}

		if user.Hash != c.UserHash {
			http.Error(w, "Wrong user-hash values", http.StatusForbidden)
			return
		}
		if user.Email == randomUserEmail {
			users, err := mysql.Misc().Users(ctx)
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
		contact, err := cs.CreateOrGet(r.Context(), c.ToContactRequest(user.ID))
		if err != nil {
			log.Println("Can't create contact error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		lead, err := orders.CreateLead(models.CreateLeadReq(c), contact)
		if err != nil {
			log.Println("Can't create lead error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		if err = orders.CreateTask(models.CreateLeadReq(c), lead); err != nil {
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
