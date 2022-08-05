package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/models"
	"github.com/hromov/jevelina/http/rest/auth"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	c := cdb.Contacts()
	var contact *models.Contact

	switch r.Method {
	case "GET":
		contact, err = c.ByID(ID)
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
		b, err := json.Marshal(contact)
		if err != nil {
			log.Println("Can't json.Marshal(contact) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(b))
	case "PUT":
		if err = json.NewDecoder(r.Body).Decode(&contact); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if uint64(contact.ID) != ID {
			http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", ID, contact.ID), http.StatusBadRequest)
			return
		}

		if err = c.DB.Omit(clause.Associations).Save(contact).Error; err != nil {
			log.Printf("Can't update contact with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		return
	case "DELETE":
		if err = c.Delete(ID); err != nil {
			log.Printf("Can't delete contact with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		return
	}

}

func ContactsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/contacts" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		contact := new(models.Contact)
		if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c := cdb.GetDB()

		user, err := auth.GetCurrentUser(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}
		contact.CreatedID = &user.ID

		if err := c.DB.Omit(clause.Associations).Create(contact).Error; err != nil {
			log.Printf("Can't create contact. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		//it actually was created ......
		b, err := json.Marshal(contact)
		if err != nil {
			log.Println("Can't json.Marshal(contact) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(b))
		return
	}

	c := cdb.Contacts()
	contactsResponse, err := c.List(FilterFromQuery(r.URL.Query()))
	if err != nil {
		log.Println("Can't get contacts error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	// log.Println("banks in main: ", banks)
	b, err := json.Marshal(contactsResponse.Contacts)
	if err != nil {
		log.Println("Can't json.Marshal(contatcts) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("X-Total-Count", strconv.FormatInt(contactsResponse.Total, 10))

	fmt.Fprint(w, string(b))
}
