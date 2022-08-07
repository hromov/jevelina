package rest

import (
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/http/rest/auth"
	api "github.com/hromov/jevelina/http/rest/handlers"
	"github.com/hromov/jevelina/http/rest/handlers/events_api"
	"github.com/hromov/jevelina/http/rest/handlers/files_api"
	"github.com/hromov/jevelina/http/rest/handlers/fin_api"
	"github.com/hromov/jevelina/useCases/orders"
)

func InitRouter(us users.Service, cs contacts.Service, ls leads.Service, os orders.Service) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/usercheck", auth.UserCheckHandler).Methods("GET")
	r.HandleFunc("/orders", api.Order(us, os)).Methods("POST")
	// TODO: uncoment for prod
	// r.Use(auth.UserCheck)
	r = UserRoutes(r, us, cs, ls)
	// TODO: uncoment for prod
	// r.Use(auth.AdminCheck)
	r = AdminRoutes(r, us)

	return r
}

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

func UserRoutes(r *mux.Router, us users.Service, cs contacts.Service, ls leads.Service) *mux.Router {
	r.HandleFunc("/contacts", api.Contacts(cs)).Methods("GET", "POST")
	r.HandleFunc("/contacts/{id}", api.Contact(cs)).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/leads", api.Leads(ls)).Methods("GET", "POST")
	r.HandleFunc("/leads/{id}", api.Lead(ls)).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/users", api.Users(us)).Methods("GET")
	r.HandleFunc("/users/{id}", api.User(us)).Methods("GET")
	r.HandleFunc("/sources", api.SourcesHandler).Methods("GET")
	r.HandleFunc("/sources/{id}", api.SourceHandler).Methods("GET")
	r.HandleFunc("/roles", api.Roles(us)).Methods("GET")
	r.HandleFunc("/roles/{id}", api.Role(us)).Methods("GET")
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
