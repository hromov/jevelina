package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/domain/misc"
	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/http/rest/auth"
)

type lead struct {
	ID        uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	ClosedAt  *time.Time
	DeletedAt *time.Time
	Name      string
	Budget    uint32
	Profit    int32

	Contact     contact
	Responsible users.User
	Created     users.User
	Step        leads.Step

	Product      misc.Product
	Manufacturer misc.Manufacturer
	Source       misc.Source
	Analytics    misc.Analytics
}

func leadFromDomain(l leads.Lead) lead {
	return lead{
		ID:        l.ID,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
		ClosedAt:  timeOrNull(l.ClosedAt),
		DeletedAt: timeOrNull(l.DeletedAt),
		Name:      l.Name,
		Budget:    l.Budget,
		Profit:    l.Profit,

		Contact:     contactFromDomain(l.Contact),
		Responsible: l.Responsible,
		Created:     l.Created,
		Step:        l.Step,

		Product:      l.Product,
		Manufacturer: l.Manufacturer,
		Source:       l.Source,
		Analytics:    l.Analytics,
	}
}

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
			_ = json.NewEncoder(w).Encode(leadFromDomain(lead))
			return
		case "PUT":
			lead := leads.LeadData{}
			if err = json.NewDecoder(r.Body).Decode(&lead); err != nil {
				log.Println("Lead decode error: ", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if uint64(lead.ID) != id {
				http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", id, lead.ID), http.StatusBadRequest)
				return
			}

			err := ls.Update(r.Context(), lead)
			if err != nil {
				log.Println("Can't save lead error: ", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				return
			}
			updated, err := ls.Get(r.Context(), id)
			if err != nil {
				log.Println("Can't get saved lead error: ", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				return
			}
			_ = json.NewEncoder(w).Encode(leadFromDomain(updated))
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
			lead := leads.LeadData{}
			if err := json.NewDecoder(r.Body).Decode(&lead); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			userValue := r.Context().Value(auth.KeyUser{})
			user, ok := userValue.(users.User)
			if !ok {
				http.Error(w, "Not a user", http.StatusForbidden)
				return
			}
			lead.ResponsibleID = user.ID
			lead.CreatedID = user.ID

			createdLead, err := ls.Create(r.Context(), lead)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				return
			}

			_ = json.NewEncoder(w).Encode(leadFromDomain(createdLead))
			return
		}

		filter, err := parseFilter(r.URL.Query())
		if err != nil {
			log.Println("Can't convert filter: ", err.Error())
			http.Error(w, "Filter error", http.StatusBadRequest)
			return
		}
		leadsResponse, err := ls.List(r.Context(), filter.toLeads())
		if err != nil {
			log.Println("Can't get leads error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		leads := make([]lead, len(leadsResponse.Leads))
		for i, l := range leadsResponse.Leads {
			leads[i] = leadFromDomain(l)
		}

		w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
		w.Header().Set("X-Total-Count", strconv.FormatInt(leadsResponse.Total, 10))
		_ = json.NewEncoder(w).Encode(leads)
	}
}
