package routes

import (
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/api/files_api"
	"github.com/hromov/jevelina/auth"
)

func FilesRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/files", auth.UserCheck(files_api.FilesHandler)).Methods("POST", "GET")
	r.HandleFunc("/files/{id}", auth.UserCheck(files_api.FileHandler)).Methods("GET")
	r.HandleFunc("/files/{id}", auth.AdminCheck(files_api.FileHandler)).Methods("DELETE")
	return r
}
