package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/http/rest/auth"
	"github.com/hromov/jevelina/useCases/tasks"
	"gorm.io/gorm"
)

func Task(ts tasks.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := getID(r)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}

		switch r.Method {
		case "GET":
			task, err := ts.Get(r.Context(), id)
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
			_ = json.NewEncoder(w).Encode(task)
			return
		case "PUT":
			task := tasks.TaskData{}
			if err = json.NewDecoder(r.Body).Decode(&task); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if uint64(task.ID) != id {
				http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", id, task.ID), http.StatusBadRequest)
				return
			}
			if task.ParentID == 0 {
				http.Error(w, "task shoud have ParentID", http.StatusBadRequest)
				return
			}

			// TODO: move all this update fileds to events
			// user, err := auth.GetCurrentUser(r)
			// if err != nil {
			// 	http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			// }
			// task.UpdatedID = &user.ID

			if err := ts.Update(r.Context(), task); err != nil {
				log.Printf("Can't update task with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		case "DELETE":

			if err := ts.Delete(r.Context(), id); err != nil {
				log.Printf("Can't delete task with ID = %d. Error: %s", id, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
	}
}

func Tasks(ts tasks.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" {
			task := tasks.TaskData{}
			if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if task.ParentID == 0 {
				http.Error(w, "task shoud have ParentID", http.StatusBadRequest)
				return
			}

			userValue := r.Context().Value(auth.KeyUser{})
			user, ok := userValue.(users.User)
			if !ok {
				http.Error(w, "Not a user", http.StatusForbidden)
				return
			}
			task.CreatedID = user.ID

			createdTask, err := ts.Create(r.Context(), task)
			if err != nil {
				log.Printf("Can't create task. Error: %s", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			_ = json.NewEncoder(w).Encode(createdTask)
			return
		}

		filter := TasksFilter(r.URL.Query())
		tasksResponse, err := ts.List(r.Context(), filter)
		if err != nil {
			log.Println("tasks Response error: ", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
		w.Header().Set("X-Total-Count", strconv.FormatInt(tasksResponse.Total, 10))
		_ = json.NewEncoder(w).Encode(tasksResponse.Tasks)
	}
}
