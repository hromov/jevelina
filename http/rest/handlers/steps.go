package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/hromov/jevelina/domain/leads"
	"gorm.io/gorm"
)

func Step(ls leads.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			step, err := ls.GetStep(r.Context(), uint8(id))
			if err != nil {
				log.Println("Can't get step error: " + err.Error())
				if errors.Is(err, gorm.ErrRecordNotFound) {
					http.NotFound(w, r)
				} else {
					http.Error(w, http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError)
				}
				return
			}
			json.NewEncoder(w).Encode(step)
			return
		case "PUT":
			step := leads.Step{}
			if err = json.NewDecoder(r.Body).Decode(&step); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if uint64(step.ID) != id {
				http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", id, step.ID), http.StatusBadRequest)
				return
			}

			if err = ls.UpdateStep(r.Context(), step); err != nil {
				log.Printf("Can't update step with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		case "DELETE":
			if err = ls.DeleteStep(r.Context(), uint8(id)); err != nil {
				log.Printf("Can't delete step with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
	}
}

func Steps(ls leads.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {
			step := leads.Step{}
			if err := json.NewDecoder(r.Body).Decode(&step); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			step, err := ls.CreateStep(r.Context(), step)
			if err != nil {
				log.Println("Can't create new step error: ", err.Error())
				http.Error(w, "Can't create step", http.StatusInternalServerError)
				return
			}

			json.NewEncoder(w).Encode(step)
			return
		}

		steps, err := ls.GetSteps(r.Context())
		if err != nil {
			log.Println("Can't get steps error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(steps)
	}
}
