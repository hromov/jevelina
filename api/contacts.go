package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hromov/cdb"
	"github.com/hromov/jevelina/base"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const contactsPageSize = 50

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	c := base.Contacts()
	var contact *cdb.Contact

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
			log.Println("Can't json.Marchal(contact) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(b))
	case "PUT":
		if err = json.NewDecoder(r.Body).Decode(&contact); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//channge to base.DB?
		if err = c.DB.Omit(clause.Associations).Save(contact).Error; err != nil {
			log.Printf("Can't update contact with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	case "DELETE":

		if err = c.DB.Delete(&cdb.Contact{ID: uint(ID)}).Error; err != nil {
			log.Printf("Can't delete contact with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	}

}

func ContactsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/contacts" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		contact := new(cdb.Contact)
		if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c := base.GetDB()
		//channge to base.DB?
		if err := c.DB.Omit(clause.Associations).Create(contact).Error; err != nil {
			log.Printf("Can't create contact. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		//it actually was created ......
		b, err := json.Marshal(contact)
		if err != nil {
			log.Println("Can't json.Marchal(contact) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(b))
		// it said that its already ok now
		// w.WriteHeader(http.StatusOK)
		return
	}

	page := r.URL.Query().Get("page")
	limit, offset := contactsPageSize, 0
	if page != "" {
		p, err := strconv.Atoi(page)
		if err == nil {
			limit = contactsPageSize
			offset = p * limit
		}
	}
	query := r.URL.Query().Get("query")
	c := base.Contacts()
	contactsResponse, err := c.List(limit, offset, query)
	if err != nil {
		log.Println("Can't get contacts error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	// log.Println("banks in main: ", banks)
	b, err := json.Marshal(contactsResponse.Contacts)
	if err != nil {
		log.Println("Can't json.Marchal(contatcts) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("X-Total-Count", strconv.FormatInt(contactsResponse.Total, 10))
	fmt.Fprintf(w, string(b))
}
