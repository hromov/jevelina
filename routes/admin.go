package routes

import (
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/api"
	"github.com/hromov/jevelina/api/events_api"
	"github.com/hromov/jevelina/api/files_api"
	"github.com/hromov/jevelina/api/fin_api"
	"github.com/hromov/jevelina/domain/users"
)

func AdminRoutes(r *mux.Router, us users.Service) *mux.Router {
	r.HandleFunc("/users", api.CreateUser(us)).Methods("POST")
	r.HandleFunc("/users/{id}", api.UpdateUser(us)).Methods("PUT")
	r.HandleFunc("/users/{id}", api.DeleteUser(us)).Methods("DELETE")
	r.HandleFunc("/sources", api.SourcesHandler).Methods("POST")
	r.HandleFunc("/sources/{id}", api.SourceHandler).Methods("PUT", "DELETE")
	r.HandleFunc("/roles", api.CreateRole(us)).Methods("POST")
	r.HandleFunc("/roles/{id}", api.UpdateRole(us)).Methods("PUT")
	r.HandleFunc("/roles/{id}", api.DeleteRole(us)).Methods("DELETE")
	r.HandleFunc("/steps", api.StepsHandler).Methods("POST")
	r.HandleFunc("/steps/{id}", api.StepHandler).Methods("PUT", "DELETE")
	r.HandleFunc("/products", api.ProductsHandler).Methods("POST")
	r.HandleFunc("/products/{id}", api.ProductHandler).Methods("PUT", "DELETE")
	r.HandleFunc("/manufacturers", api.ManufacturersHandler).Methods("POST")
	r.HandleFunc("/manufacturers/{id}", api.ManufacturerHandler).Methods("PUT", "DELETE")
	r.HandleFunc("/tags", api.TagsHandler).Methods("POST")
	r.HandleFunc("/tags/{id}", api.TagHandler).Methods("PUT", "DELETE")
	r.HandleFunc("/tasks/{id}", api.TaskHandler).Methods("DELETE")
	r.HandleFunc("/tasktypes", api.TaskTypesHandler).Methods("POST")
	r.HandleFunc("/tasktypes/{id}", api.TaskTypeHandler).Methods("PUT", "DELETE")
	r.HandleFunc("/events/transfers", events_api.ListHandler).Methods("GET")
	r.HandleFunc("/files/{id}", files_api.FileHandler).Methods("DELETE")
	// Finance part
	r.HandleFunc("/wallets", fin_api.WalletsHandler).Methods("POST")
	r.HandleFunc("/wallets/{id}", fin_api.WalletHandler).Methods("PUT", "DELETE")
	r.HandleFunc("/wallets/{id}/close", fin_api.CloseWalletHandler).Methods("GET")
	r.HandleFunc("/wallets/{id}/open", fin_api.OpenWalletHandler).Methods("GET")
	r.HandleFunc("/transfers/{id}", fin_api.TransferHandler).Methods("DELETE")
	r.HandleFunc("/transfers/{id}/complete", fin_api.CompleteTransferHandler).Methods("GET")
	return r
}
