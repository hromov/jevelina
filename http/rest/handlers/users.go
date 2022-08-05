package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/domain/users"
	"gorm.io/gorm"
)

type user struct {
	ID           uint64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
	Name         string
	Email        string
	Hash         string
	Distribution float32
	Role         string
}

type changeUser struct {
	ID           uint64
	Name         string
	Email        string
	Hash         string
	Distribution float32
	RoleID       uint8
}

func fromUser(u users.User) user {
	return user(u)
}

func (u changeUser) toDomain() users.ChangeUser {
	return users.ChangeUser(u)
}

func CreateUser(us users.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		newUser := changeUser{}
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			log.Println("user decode error: ", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		user, err := us.Create(r.Context(), newUser.toDomain())
		if err != nil {
			log.Println("create user error: ", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
	}
}

func UpdateUser(us users.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}

		newUser := changeUser{}
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			log.Println("user decode error: ", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if id != newUser.ID {
			http.Error(w, "Route user ID and user.ID doesn't match", http.StatusBadRequest)
			return
		}

		if err := us.Update(r.Context(), newUser.toDomain()); err != nil {
			log.Println("update user error: ", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteUser(us users.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := us.Delete(r.Context(), id); err != nil {
			log.Println("delete user error: ", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func User(us users.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.ParseUint(vars["id"], 10, 32)
		if err != nil {
			http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
			return
		}
		user, err := us.Get(r.Context(), id)
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
	}
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
