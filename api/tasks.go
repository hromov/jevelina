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
		task, err = c.Task(ID)
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
		if task.ParentID == 0 {
			http.Error(w, "task shoud have ParentID", http.StatusBadRequest)
			return
		}
		//TODO: change after AUTH!
		task.UpdatedID = task.ResponsibleID

		//channge to base.DB?
		if err = c.DB.Save(task).Error; err != nil {
			log.Printf("Can't update task with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	case "DELETE":

		if err = c.DB.Delete(&models.Task{ID: ID}).Error; err != nil {
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
		if task.ParentID == 0 {
			http.Error(w, "task shoud have ParentID", http.StatusBadRequest)
			return
		}
		//TODO: CHANGE TO REAL, AFTER AUTH!!!!!!!!!!!!!!!!!!!!!!!!!!!
		task.CreatedID = task.ResponsibleID
		c := base.GetDB()
		//channge to base.DB?
		if err := c.DB.Omit(clause.Associations).Create(task).Error; err != nil {
			log.Printf("Can't create task. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		//fix after 64 update
		fullTask, err := c.Misc().Task(task.ID)
		if err != nil {
			log.Printf("Task should be created but we wasn't able to get it back. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		//it actually was created ......
		b, err := json.Marshal(fullTask)
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

	tasksResponse, err := c.Tasks(filterFromQuery(r.URL.Query()))
	// log.Println("banks in main: ", banks)
	b, err := json.Marshal(tasksResponse.Tasks)
	if err != nil {
		log.Println("Can't json.Marshal(tasks) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("X-Total-Count", strconv.FormatInt(tasksResponse.Total, 10))
	fmt.Fprintf(w, string(b))
}
