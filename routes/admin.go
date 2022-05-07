package routes

import (
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/api"
	"github.com/hromov/jevelina/auth"
)

func AdminRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/users", auth.AdminCheck(api.UsersHandler)).Methods("POST")
	r.HandleFunc("/users/{id}", auth.AdminCheck(api.UserHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/sources", auth.AdminCheck(api.SourcesHandler)).Methods("POST")
	r.HandleFunc("/sources/{id}", auth.AdminCheck(api.SourceHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/roles", auth.AdminCheck(api.RolesHandler)).Methods("POST")
	r.HandleFunc("/roles/{id}", auth.AdminCheck(api.RoleHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/steps", auth.AdminCheck(api.StepsHandler)).Methods("POST")
	r.HandleFunc("/steps/{id}", auth.AdminCheck(api.StepHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/products", auth.AdminCheck(api.ProductsHandler)).Methods("POST")
	r.HandleFunc("/products/{id}", auth.AdminCheck(api.ProductHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/manufacturers", auth.AdminCheck(api.ManufacturersHandler)).Methods("POST")
	r.HandleFunc("/manufacturers/{id}", auth.AdminCheck(api.ManufacturerHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/tags", auth.AdminCheck(api.TagsHandler)).Methods("POST")
	r.HandleFunc("/tags/{id}", auth.AdminCheck(api.TagHandler)).Methods("PUT", "DELETE")
	r.HandleFunc("/tasks/{id}", auth.AdminCheck(api.TaskHandler)).Methods("DELETE")
	r.HandleFunc("/tasktypes", auth.AdminCheck(api.TaskTypesHandler)).Methods("POST")
	r.HandleFunc("/tasktypes/{id}", auth.AdminCheck(api.TaskTypeHandler)).Methods("PUT", "DELETE")
	return r
}
