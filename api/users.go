package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/cdb"
	"github.com/hromov/jevelina/cdb/models"
	"github.com/hromov/jevelina/domain/users"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type user struct {
	ID           uint64    `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Hash         string    `json:"hash"`
	Distribution float32   `json:"distribution"`
	Role         string    `json:"role"`
}

func fromUser(u users.User) user {
	return user(u)
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	c := cdb.Misc()
	var user users.User

	switch r.Method {
	case "GET":
		user, err = c.User(r.Context(), ID)
		if err != nil {
			log.Println("Can't get user error: " + err.Error())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.NotFound(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
		b, err := json.Marshal(user)
		if err != nil {
			log.Println("Can't json.Marshal(user) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(b))
	case "PUT":
		if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if uint64(user.ID) != ID {
			http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", ID, user.ID), http.StatusBadRequest)
			return
		}

		//channge to base.DB?
		if err = c.DB.Omit(clause.Associations).Save(user).Error; err != nil {
			log.Printf("Can't update user with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	case "DELETE":
		if err = c.DB.Delete(&models.User{ID: ID}).Error; err != nil {
			log.Printf("Can't delete user with ID = %d. Error: %s", ID, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		// w.WriteHeader(http.StatusOK)
		return
	}

}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		user := new(models.User)
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		c := cdb.Misc()
		//channge to base.DB?
		if err := c.DB.Omit(clause.Associations).Create(user).Error; err != nil {
			log.Printf("Can't create user. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
		//remove uint conversion when cdb updated
		fullUser, err := c.User(r.Context(), user.ID)
		if err != nil {
			log.Printf("User should be created but we wasn't able to get it back. Error: %s", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		//it actually was created ......
		b, err := json.Marshal(fullUser)
		if err != nil {
			log.Println("Can't json.Marshal(user) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(b))
		return
	}

	// c := cdb.Misc()
	// usersResponse, err := c.Users()
	// if err != nil {
	// 	log.Println("Can't get users error: " + err.Error())
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError),
	// 		http.StatusInternalServerError)
	// }
	// // log.Println("banks in main: ", banks)
	// b, err := json.Marshal(usersResponse)
	// if err != nil {
	// 	log.Println("Can't json.Marshal(contatcts) error: " + err.Error())
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError),
	// 		http.StatusInternalServerError)
	// 	return
	// }
	// total := strconv.Itoa(len(usersResponse))
	// w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	// w.Header().Set("X-Total-Count", total)
	// fmt.Fprint(w, string(b))
}

func Users(us users.Service) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		users, err := us.List(r.Context())
		if err != nil {
			log.Println("Can't get users error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}

		usersResponse := make([]user, len(users))
		for i, u := range users {
			usersResponse[i] = fromUser(u)
		}

		b, err := json.Marshal(usersResponse)
		if err != nil {
			log.Println("Can't json.Marshal(users) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		total := strconv.Itoa(len(usersResponse))
		w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
		w.Header().Set("X-Total-Count", total)
		fmt.Fprint(w, string(b))
	}
}
