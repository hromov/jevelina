package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/hromov/cdb/leads"
	"github.com/hromov/jevelina/base"
)

const leadsPageSize = 50

func LeadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}
	l := base.Leads()
	lead, err := l.ByID(ID)
	if err != nil {
		log.Println("Can't get lead error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	b, err := json.Marshal(lead)
	if err != nil {
		log.Println("Can't json.Marchal(contatct) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(b))
}

func LeadsHandler(w http.ResponseWriter, r *http.Request) {
	leadsResponse := &leads.LeadsResponse{}
	var err error

	if r.URL.Path != "/leads" {
		http.NotFound(w, r)
		return
	}

	page := r.URL.Query().Get("page")
	limit, offset := leadsPageSize, 0
	if page != "" {
		p, err := strconv.Atoi(page)
		if err == nil {
			limit = leadsPageSize
			offset = p * limit
		}
	}
	query := r.URL.Query().Get("query")
	contactID := r.URL.Query().Get("contactID")
	l := base.Leads()
	if contactID != "" {
		ID, err := strconv.ParseUint(contactID, 10, 32)
		if err != nil {
			http.Error(w, "clientID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}
		leadsResponse, err = l.ByContact(uint(ID))
		if err != nil {
			log.Println("Can't get leads error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	} else {
		leadsResponse, err = l.List(limit, offset, query)
		if err != nil {
			log.Println("Can't get leads error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}

	// log.Println("banks in main: ", banks)
	b, err := json.Marshal(leadsResponse.Leads)
	if err != nil {
		log.Println("Can't json.Marchal(contatcts) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("X-Total-Count", strconv.FormatInt(leadsResponse.Total, 10))
	fmt.Fprintf(w, string(b))
}
