package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/storage/mysql"
	"gorm.io/gorm"
)

type createRoleRequest struct {
	Priority uint8  `validate:"required"`
	Role     string `validate:"required"`
}

type updateRoleRequest struct {
	ID       uint8  `validate:"required"`
	Priority uint8  `validate:"required"`
	Role     string `validate:"required"`
}

func (c *createRoleRequest) toDomain() users.Role {
	return users.Role{
		Priority: c.Priority,
		Role:     c.Role,
	}
}

func (role *updateRoleRequest) toDomain() users.Role {
	return users.Role{
		ID:       role.ID,
		Priority: role.Priority,
		Role:     role.Role,
	}
}

func RoleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "ID conversion error: "+err.Error(), http.StatusBadRequest)
		return
	}

	c := mysql.Misc()
	var role users.Role

	switch r.Method {
	case "GET":
		role, err = c.Role(r.Context(), uint8(ID))
		if err != nil {
			log.Println("Can't get role error: " + err.Error())
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.NotFound(w, r)
			} else {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
			return
		}
		b, err := json.Marshal(role)
		if err != nil {
			log.Println("Can't json.Marshal(role) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(b))
	}

}

func RolesHandler(w http.ResponseWriter, r *http.Request) {
	c := mysql.Misc()
	rolesResponse, err := c.Roles(r.Context())
	if err != nil {
		log.Println("Can't get roles error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	}
	// log.Println("banks in main: ", banks)
	b, err := json.Marshal(rolesResponse)
	if err != nil {
		log.Println("Can't json.Marshal(contatcts) error: " + err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
		return
	}
	total := strconv.Itoa(len(rolesResponse))
	w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count")
	w.Header().Set("X-Total-Count", total)
	fmt.Fprint(w, string(b))
}

func CreateRole(us users.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		createRequest := createRoleRequest{}
		if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
			log.Println("create role decode error: ", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		role, err := us.CreateRole(r.Context(), createRequest.toDomain())
		if err != nil {
			log.Println("create user error: ", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(role)
		if err != nil {
			log.Println("Can't json.Marshal(user) error: " + err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		log.Println("new role was created: ", string(b))
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, string(b))
	}
}

func UpdateRole(us users.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ID, err := getRouteID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updateRequest := updateRoleRequest{}
		if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
			log.Println("create role decode error: ", err.Error())
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if updateRequest.ID != uint8(ID) {
			http.Error(w, fmt.Sprintf("url ID = %d is not the one from the request: %d", ID, updateRequest.ID), http.StatusBadRequest)
			return
		}

		if err := us.UpdateRole(r.Context(), updateRequest.toDomain()); err != nil {
			log.Println("can't update user role error: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteRole(us users.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ID, err := getRouteID(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := us.DeleteRole(r.Context(), uint8(ID)); err != nil {
			log.Println("Can't delete role error: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
