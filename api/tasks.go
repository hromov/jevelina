package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/auth"
	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/models"
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

	c := cdb.Misc()
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
			log.Println("Can't json.Marshal(task) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, b)
		return
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

		user, err := auth.GetCurrentUser(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}
		task.UpdatedID = &user.ID

		if err = c.DB.Omit(clause.Associations).Save(task).Error; err != nil {
			log.Printf("Can't update task with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		return
	case "DELETE":

		if err = c.DB.Delete(&models.Task{ID: ID}).Error; err != nil {
			log.Printf("Can't delete task with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		return
	}

}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/tasks" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		task := new(models.Task)
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if task.ParentID == 0 {
			http.Error(w, "task shoud have ParentID", http.StatusBadRequest)
			return
		}

		user, err := auth.GetCurrentUser(r)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		}
		task.CreatedID = &user.ID

		c := cdb.GetDB()
		if err := c.DB.Omit(clause.Associations).Create(task).Error; err != nil {
			log.Printf("Can't create task. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		fullTask, err := cdb.Misc().Task(task.ID)
		if err != nil {
			log.Printf("Task should be created but we wasn't able to get it back. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		//it actually was created ......
		b, err := json.Marshal(fullTask)
		if err != nil {
			log.Println("Can't json.Marshal(task) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, b)
		return
	}

	c := cdb.Misc()

	tasksResponse, err := c.Tasks(FilterFromQuery(r.URL.Query()))
	if err != nil {
		log.Println("tasks Response error: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
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
	fmt.Fprint(w, b)
}
