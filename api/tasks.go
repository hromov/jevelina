package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hromov/cdb/models"
	"github.com/hromov/jevelina/base"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	c := base.GetDB().Misc()
	var task *models.Task

	switch r.Method {
	case "GET":
		task, err = c.Task(uint(ID))
		if err != nil {
			log.Println("Can't get task error: " + err.Error())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.NotFound(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
		b, err := json.Marshal(task)
		if err != nil {
			log.Println("Can't json.Marchal(task) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(b))
	case "PUT":
		if err = json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if uint64(task.ID) != ID {
			http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", ID, task.ID), http.StatusBadRequest)
			return
		}
		if task.LeadID == 0 && task.ContactID == 0 {
			http.Error(w, "task shoud have LeadID or ContactID", http.StatusBadRequest)
			return
		}
		//channge to base.DB?
		if err = c.DB.Save(task).Error; err != nil {
			log.Printf("Can't update task with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	case "DELETE":

		if err = c.DB.Delete(&models.Task{ID: uint(ID)}).Error; err != nil {
			log.Printf("Can't delete task with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	}

}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tasks" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		// CHECK ID for LEAD OR CONTACT
		task := new(models.Task)
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if task.LeadID == 0 && task.ContactID == 0 {
			http.Error(w, "task shoud have LeadID or ContactID", http.StatusBadRequest)
			return
		}
		c := base.GetDB()
		//channge to base.DB?
		if err := c.DB.Omit(clause.Associations).Create(task).Error; err != nil {
			log.Printf("Can't create task. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		//it actually was created ......
		b, err := json.Marshal(task)
		if err != nil {
			log.Println("Can't json.Marchal(task) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(b))
		// it said that its already ok now
		// w.WriteHeader(http.StatusOK)
		return
	}

	c := base.GetDB().Misc()

	contactID := r.URL.Query().Get("contactID")
	leadID := r.URL.Query().Get("leadID")

	if contactID == "" && leadID == "" {
		http.Error(w, "ContactID or LeadID is required to get tasks", http.StatusBadRequest)
		return
	}

	tasks := []models.Task{}
	if contactID != "" {
		ID, err := strconv.ParseUint(contactID, 10, 64)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}
		tasks, err = c.TasksByContact(uint(ID))
	} else {
		ID, err := strconv.ParseUint(contactID, 10, 64)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}
		tasks, err = c.TasksByLead(uint(ID))
	}
	// log.Println("banks in main: ", banks)
	b, err := json.Marshal(tasks)
	if err != nil {
		log.Println("Can't json.Marchal(tasks) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	total := strconv.Itoa(len(tasks))
	w.Header().Set("X-Total-Count", total)
	fmt.Fprintf(w, string(b))
}
