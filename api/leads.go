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
	"github.com/hromov/jevelina/base"
	"github.com/hromov/jevelina/cdb/models"
)

func LeadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	l := base.GetDB().Leads()
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
		fmt.Fprintf(w, string(b))
	case "PUT":
		if err = json.NewDecoder(r.Body).Decode(&lead); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//channge to base.DB?
		if uint64(lead.ID) != ID {
			http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", ID, lead.ID), http.StatusBadRequest)
			return
		}
		if err = l.DB.Save(lead).Error; err != nil {
			log.Printf("Can't update lead with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	case "DELETE":

		if err = l.DB.Delete(&models.Lead{ID: ID}).Error; err != nil {
			log.Printf("Can't delete lead with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	}
}

func LeadsHandler(w http.ResponseWriter, r *http.Request) {
	leadsResponse := &models.LeadsResponse{}
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
		c := base.GetDB()

		if err := c.DB.Create(lead).Error; err != nil {
			log.Printf("Can't create lead. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		//it actually was created ......
		b, err := json.Marshal(lead)
		if err != nil {
			log.Println("Can't json.Marchal(lead) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(b))
		return
	}

	email, err := auth.GetMailByToken(r)
	if err != nil {
		log.Println("Error recieving email: ", err.Error())
	}
	log.Println("Email = ", email)

	l := base.GetDB().Leads()
	leadsResponse, err = l.List(filterFromQuery(r.URL.Query()))
	if err != nil {
		log.Println("Can't get leads error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}

	// log.Println("banks in main: ", banks)
	b, err := json.Marshal(leadsResponse.Leads)
	if err != nil {
		log.Println("Can't json.Marchal(contatcts) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("X-Total-Count", strconv.FormatInt(leadsResponse.Total, 10))
	fmt.Fprintf(w, string(b))
}
