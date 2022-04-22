package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hromov/cdb"
)

const contactsPageSize = 50

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}
	contact, err := cdb.ContactByID(ID)
	if err != nil {
		log.Println("Can't get contact error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	b, err := json.Marshal(contact)
	if err != nil {
		log.Println("Can't json.Marchal(contatct) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(b))
}

func ContactsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/contacts" {
		http.NotFound(w, r)
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

	contactsResponse, err := cdb.Contacts(limit, offset, query)
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
