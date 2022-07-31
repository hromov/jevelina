package routes

import (
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/api"
	"github.com/hromov/jevelina/api/files_api"
	"github.com/hromov/jevelina/api/fin_api"
	"github.com/hromov/jevelina/domain/users"
)

func UserRoutes(r *mux.Router, us users.Service) *mux.Router {
	r.HandleFunc("/contacts", api.ContactsHandler).Methods("GET", "POST")
	r.HandleFunc("/contacts/{id}", api.ContactHandler).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/leads", api.LeadsHandler).Methods("GET", "POST")
	r.HandleFunc("/leads/{id}", api.LeadHandler).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/users", api.Users(us)).Methods("GET")
	r.HandleFunc("/users/{id}", api.User(us)).Methods("GET")
	r.HandleFunc("/sources", api.SourcesHandler).Methods("GET")
	r.HandleFunc("/sources/{id}", api.SourceHandler).Methods("GET")
	r.HandleFunc("/roles", api.RolesHandler).Methods("GET")
	r.HandleFunc("/roles/{id}", api.RoleHandler).Methods("GET")
	r.HandleFunc("/steps", api.StepsHandler).Methods("GET")
	r.HandleFunc("/steps/{id}", api.StepHandler).Methods("GET")
	r.HandleFunc("/products", api.ProductsHandler).Methods("GET")
	r.HandleFunc("/products/{id}", api.ProductHandler).Methods("GET")
	r.HandleFunc("/manufacturers", api.ManufacturersHandler).Methods("GET")
	r.HandleFunc("/manufacturers/{id}", api.ManufacturerHandler).Methods("GET")
	r.HandleFunc("/tags", api.TagsHandler).Methods("GET")
	r.HandleFunc("/tags/{id}", api.TagHandler).Methods("GET")
	r.HandleFunc("/tasks", api.TasksHandler).Methods("GET", "POST")
	r.HandleFunc("/tasks/{id}", api.TaskHandler).Methods("GET", "PUT")
	r.HandleFunc("/tasktypes", api.TaskTypesHandler).Methods("GET")
	r.HandleFunc("/tasktypes/{id}", api.TaskTypeHandler).Methods("GET")
	r.HandleFunc("/files", files_api.FilesHandler).Methods("POST", "GET")
	r.HandleFunc("/files/{id}", files_api.FileHandler).Methods("GET")

	r.HandleFunc("/wallets", fin_api.WalletsHandler).Methods("GET")
	r.HandleFunc("/transfers", fin_api.TransfersHandler).Methods("GET", "POST")
	r.HandleFunc("/transfers/{id}", fin_api.TransferHandler).Methods("PUT")
	r.HandleFunc("/categories", fin_api.CategoriesHandler).Methods("GET")
	r.HandleFunc("/analytics/categories", fin_api.CategoriesSumHandler).Methods("GET")
	return r
}
