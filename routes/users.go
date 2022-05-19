package routes

import (
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/api"
	"github.com/hromov/jevelina/auth"
)

func UserRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/contacts", auth.UserCheck(api.ContactsHandler)).Methods("GET", "POST")
	r.HandleFunc("/contacts/{id}", auth.UserCheck(api.ContactHandler)).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/leads", auth.UserCheck(api.LeadsHandler)).Methods("GET", "POST")
	r.HandleFunc("/leads/{id}", auth.UserCheck(api.LeadHandler)).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/users", auth.UserCheck(api.UsersHandler)).Methods("GET")
	r.HandleFunc("/users/{id}", auth.UserCheck(api.UserHandler)).Methods("GET")
	r.HandleFunc("/sources", auth.UserCheck(api.SourcesHandler)).Methods("GET")
	r.HandleFunc("/sources/{id}", auth.UserCheck(api.SourceHandler)).Methods("GET")
	r.HandleFunc("/roles", auth.UserCheck(api.RolesHandler)).Methods("GET")
	r.HandleFunc("/roles/{id}", auth.UserCheck(api.RoleHandler)).Methods("GET")
	r.HandleFunc("/steps", auth.UserCheck(api.StepsHandler)).Methods("GET")
	r.HandleFunc("/steps/{id}", auth.UserCheck(api.StepHandler)).Methods("GET")
	r.HandleFunc("/products", auth.UserCheck(api.ProductsHandler)).Methods("GET")
	r.HandleFunc("/products/{id}", auth.UserCheck(api.ProductHandler)).Methods("GET")
	r.HandleFunc("/manufacturers", auth.UserCheck(api.ManufacturersHandler)).Methods("GET")
	r.HandleFunc("/manufacturers/{id}", auth.UserCheck(api.ManufacturerHandler)).Methods("GET")
	r.HandleFunc("/tags", auth.UserCheck(api.TagsHandler)).Methods("GET")
	r.HandleFunc("/tags/{id}", auth.UserCheck(api.TagHandler)).Methods("GET")
	r.HandleFunc("/tasks", auth.UserCheck(api.TasksHandler)).Methods("GET", "POST")
	r.HandleFunc("/tasks/{id}", auth.UserCheck(api.TaskHandler)).Methods("GET", "PUT")
	r.HandleFunc("/tasktypes", auth.UserCheck(api.TaskTypesHandler)).Methods("GET")
	r.HandleFunc("/tasktypes/{id}", auth.UserCheck(api.TaskTypeHandler)).Methods("GET")
	return r
}
