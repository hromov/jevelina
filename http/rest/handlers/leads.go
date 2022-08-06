package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/http/rest/auth"
)

func Lead(ls leads.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			lead, err := ls.Get(r.Context(), id)
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
			_ = json.NewEncoder(w).Encode(lead)
			return
		case "PUT":
			lead := leads.Lead{}
			if err = json.NewDecoder(r.Body).Decode(&lead); err != nil {
				log.Println("Lead decode error: ", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if uint64(lead.ID) != id {
				http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", id, lead.ID), http.StatusBadRequest)
				return
			}

			if err := ls.Update(r.Context(), lead); err != nil {
				log.Println("Can't save lead error: ", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		case "DELETE":
			if err := ls.Delete(r.Context(), id); err != nil {
				log.Printf("Can't delete lead with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}

func Leads(ls leads.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {
			lead := leads.Lead{}
			if err := json.NewDecoder(r.Body).Decode(&lead); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			user, err := auth.GetCurrentUser(r)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			}
			lead.Responsible = user
			lead.Created = user

			if lead, err = ls.Create(r.Context(), lead); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				return
			}

			_ = json.NewEncoder(w).Encode(lead)
			return
		}

		leadsResponse, err := ls.List(r.Context(), LeadsFilter(r.URL.Query()))
		if err != nil {
			log.Println("Can't get leads error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
		w.Header().Set("X-Total-Count", strconv.FormatInt(leadsResponse.Total, 10))
		_ = json.NewEncoder(w).Encode(leadsResponse.Leads)
	}
}
