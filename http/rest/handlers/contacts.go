package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/misc"
	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/http/rest/auth"
	"gorm.io/gorm"
)

type contactRequest struct {
	ID            uint64
	Name          string
	SecondName    string
	ResponsibleID uint64
	CreatedID     uint64
	Phone         string
	SecondPhone   string
	Email         string
	SecondEmail   string
	URL           string

	City    string
	Address string

	SourceID uint8
	Position string

	Analytics misc.Analytics
}

func (c *contactRequest) toDomain() contacts.ContactRequest {
	return contacts.ContactRequest{
		ID:            c.ID,
		Name:          c.Name,
		SecondName:    c.SecondName,
		ResponsibleID: c.ResponsibleID,
		CreatedID:     c.CreatedID,
		Phone:         c.Phone,
		SecondPhone:   c.SecondPhone,
		Email:         c.Email,
		SecondEmail:   c.SecondEmail,
		URL:           c.URL,

		City:    c.City,
		Address: c.Address,

		SourceID: c.SourceID,
		Position: c.Position,

		Analytics: c.Analytics,
	}
}

type contact struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Name        string
	SecondName  string
	Responsible users.User
	Created     users.User
	Phone       string
	SecondPhone string
	Email       string
	SecondEmail string
	URL         string

	City    string
	Address string

	Source   misc.Source
	Position string

	Analytics misc.Analytics
}

func contactFromDomain(c contacts.Contact) contact {
	return contact{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: timeOrNull(c.DeletedAt),

		Name:        c.Name,
		SecondName:  c.SecondName,
		Responsible: c.Responsible,
		Created:     c.Created,
		Phone:       c.Phone,
		SecondPhone: c.SecondPhone,
		Email:       c.Email,
		SecondEmail: c.SecondEmail,
		URL:         c.URL,

		City:    c.City,
		Address: c.Address,

		Source:   c.Source,
		Position: c.Position,

		Analytics: c.Analytics,
	}
}

func Contact(cs contacts.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			contact, err := cs.Get(r.Context(), id)
			if err != nil {
				log.Println("Can't get contact error: " + err.Error())
				if errors.Is(err, gorm.ErrRecordNotFound) {
					http.NotFound(w, r)
				} else {
					http.Error(w, http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError)
				}
				return
			}
			encode(w, contactFromDomain(contact))
		case "PUT":
			contact := contactRequest{}
			if err = json.NewDecoder(r.Body).Decode(&contact); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if contact.ID != id {
				http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", id, contact.ID), http.StatusBadRequest)
				return
			}

			if err := cs.Update(r.Context(), contact.toDomain()); err != nil {
				log.Printf("Can't update contact with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			updated, err := cs.Get(r.Context(), id)
			if err != nil {
				log.Println("Can't get saved contact error: ", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				return
			}
			encode(w, contactFromDomain(updated))
			return
		case "DELETE":
			if err := cs.Delete(r.Context(), id); err != nil {
				log.Printf("Can't delete contact with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
	}
}

func Contacts(cs contacts.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			contact := contactRequest{}
			if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			userValue := r.Context().Value(auth.KeyUser{})
			user, ok := userValue.(users.User)
			if !ok {
				http.Error(w, "Not a user", http.StatusForbidden)
				return
			}
			contact.CreatedID = user.ID

			createdContact, err := cs.Create(r.Context(), contact.toDomain())
			if err != nil {
				log.Printf("Can't create contact. Error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}

			encode(w, contactFromDomain(createdContact))
			return
		}
		filter, err := parseFilter(r.URL.Query())
		if err != nil {
			log.Println("Can't convert filter: ", err.Error())
			http.Error(w, "Filter error", http.StatusBadRequest)
			return
		}
		contactsResponse, err := cs.List(r.Context(), filter.toContacts())
		if err != nil {
			log.Println("Can't get contacts error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		contacts := make([]contact, len(contactsResponse.Contacts))
		for i, c := range contactsResponse.Contacts {
			contacts[i] = contactFromDomain(c)
		}

		w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
		w.Header().Set("X-Total-Count", strconv.FormatInt(contactsResponse.Total, 10))
		w.WriteHeader(http.StatusCreated)
		encode(w, contacts)
	}
}
