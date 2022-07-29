package api

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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func StepHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	c := cdb.Misc()
	var step *models.Step

	switch r.Method {
	case "GET":
		step, err = c.Step(uint8(ID))
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
		b, err := json.Marshal(step)
		if err != nil {
			log.Println("Can't json.Marshal(step) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(b)
	case "PUT":
		if err = json.NewDecoder(r.Body).Decode(&step); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if uint64(step.ID) != ID {
			http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", ID, step.ID), http.StatusBadRequest)
			return
		}

		//channge to base.DB?
		if err = c.DB.Omit(clause.Associations).Save(step).Error; err != nil {
			log.Printf("Can't update step with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	case "DELETE":

		if err = c.DB.Delete(&models.Step{ID: uint8(ID)}).Error; err != nil {
			log.Printf("Can't delete step with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	}

}

func StepsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/steps" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		step := new(models.Step)
		if err := json.NewDecoder(r.Body).Decode(&step); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c := cdb.GetDB()
		//channge to base.DB?
		if err := c.DB.Omit(clause.Associations).Create(step).Error; err != nil {
			log.Printf("Can't create step. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		//it actually was created ......
		b, err := json.Marshal(step)
		if err != nil {
			log.Println("Can't json.Marshal(step) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(b)
		// it said that its already ok now
		// w.WriteHeader(http.StatusOK)
		return
	}

	c := cdb.Misc()
	stepsResponse, err := c.Steps()
	if err != nil {
		log.Println("Can't get steps error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	// log.Println("banks in main: ", banks)
	b, err := json.Marshal(stepsResponse)
	if err != nil {
		log.Println("Can't json.Marshal(contatcts) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	total := strconv.Itoa(len(stepsResponse))
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("X-Total-Count", total)
	_, _ = w.Write(b)
}
