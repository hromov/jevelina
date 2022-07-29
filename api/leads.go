package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"github.com/hromov/jevelina/auth"
	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/models"
)

func LeadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	l := cdb.Leads()
	lead := new(models.Lead)
	switch r.Method {
	case "GET":
		lead, err = l.ByID(ID)
		if err != nil {
			log.Println("Can't get lead error: " + err.Error())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.NotFound(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
		b, err := json.Marshal(lead)
		if err != nil {
			log.Println("Can't json.Marshal(lead) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, b)
	case "PUT":
		if err = json.NewDecoder(r.Body).Decode(&lead); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if uint64(lead.ID) != ID {
			http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", ID, lead.ID), http.StatusBadRequest)
			return
		}

		if _, err := l.Save(lead); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		return
	case "DELETE":

		if err = l.Delete(ID); err != nil {
			log.Printf("Can't delete lead with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		return
	}
}

func LeadsHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.URL.Path != "/leads" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		lead := new(models.Lead)
		if err := json.NewDecoder(r.Body).Decode(&lead); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := auth.GetCurrentUser(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}
		lead.ResponsibleID = &user.ID
		lead.CreatedID = &user.ID

		if lead, err = cdb.Leads().Save(lead); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(lead)
		if err != nil {
			log.Println("Can't json.Marshal(lead) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, b)
		return
	}

	l := cdb.Leads()
	leadsResponse, err := l.List(FilterFromQuery(r.URL.Query()))
	if err != nil {
		log.Println("Can't get leads error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

	b, err := json.Marshal(leadsResponse.Leads)
	if err != nil {
		log.Println("Can't json.Marshal(contatcts) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("X-Total-Count", strconv.FormatInt(leadsResponse.Total, 10))
	fmt.Fprint(w, b)
}
