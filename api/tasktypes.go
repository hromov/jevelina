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

func TaskTypeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	c := cdb.Misc()
	var tasktype *models.TaskType

	switch r.Method {
	case "GET":
		tasktype, err = c.TaskType(uint8(ID))
		if err != nil {
			log.Println("Can't get tasktype error: " + err.Error())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.NotFound(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
		b, err := json.Marshal(tasktype)
		if err != nil {
			log.Println("Can't json.Marshal(tasktype) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(b))
	case "PUT":
		if err = json.NewDecoder(r.Body).Decode(&tasktype); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if uint64(tasktype.ID) != ID {
			http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", ID, tasktype.ID), http.StatusBadRequest)
			return
		}
		//channge to base.DB?
		if err = c.DB.Save(tasktype).Error; err != nil {
			log.Printf("Can't update tasktype with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	case "DELETE":

		if err = c.DB.Delete(&models.TaskType{ID: uint8(ID)}).Error; err != nil {
			log.Printf("Can't delete tasktype with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	}

}

func TaskTypesHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tasktypes" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		tasktype := new(models.TaskType)
		if err := json.NewDecoder(r.Body).Decode(&tasktype); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c := cdb.GetDB()
		//channge to base.DB?
		if err := c.DB.Omit(clause.Associations).Create(tasktype).Error; err != nil {
			log.Printf("Can't create tasktype. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		//it actually was created ......
		b, err := json.Marshal(tasktype)
		if err != nil {
			log.Println("Can't json.Marshal(tasktype) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(b))
		// it said that its already ok now
		// w.WriteHeader(http.StatusOK)
		return
	}

	c := cdb.Misc()
	tasktypesResponse, err := c.TaskTypes()
	if err != nil {
		log.Println("Can't get tasktypes error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	// log.Println("banks in main: ", banks)
	b, err := json.Marshal(tasktypesResponse)
	if err != nil {
		log.Println("Can't json.Marshal(contatcts) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	total := strconv.Itoa(len(tasktypesResponse))
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("X-Total-Count", total)
	fmt.Fprint(w, string(b))
}
