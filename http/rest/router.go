package rest

import (
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/domain/contacts"
	"github.com/hromov/jevelina/domain/finances"
	"github.com/hromov/jevelina/domain/leads"
	"github.com/hromov/jevelina/domain/misc"
	"github.com/hromov/jevelina/domain/misc/files"
	"github.com/hromov/jevelina/domain/users"
	"github.com/hromov/jevelina/http/rest/auth"
	api "github.com/hromov/jevelina/http/rest/handlers"
	"github.com/hromov/jevelina/http/rest/handlers/events_api"
	"github.com/hromov/jevelina/useCases/orders"
	"github.com/hromov/jevelina/useCases/tasks"
)

func InitRouter(
	us users.Service, cs contacts.Service, ls leads.Service,
	os orders.Service, ms misc.Service, ts tasks.Service,
	as auth.Service, fs files.Service, fin finances.Service,
) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/usercheck", as.UserCheckHandler()).Methods("GET")
	r.HandleFunc("/orders", api.Order(us, os)).Methods("POST")
	// TODO: uncoment for prod
	// r.Use(as.UserCheck)
	r = UserRoutes(r, us, cs, ls, ms, ts, fs, fin)
	// TODO: uncoment for prod
	// r.Use(as.AdminCheck)
	r = AdminRoutes(r, us, ms, ls, ts, fs, fin)

	return r
}

func AdminRoutes(
	r *mux.Router, us users.Service, ms misc.Service,
	ls leads.Service, ts tasks.Service, fs files.Service,
	fin finances.Service,
) *mux.Router {
	r.HandleFunc("/users", api.CreateUser(us)).Methods("POST")
	r.HandleFunc("/users/{id}", api.UpdateUser(us)).Methods("PUT")
	r.HandleFunc("/users/{id}", api.DeleteUser(us)).Methods("DELETE")
	r.HandleFunc("/sources", api.Sources(ms)).Methods("POST")
	r.HandleFunc("/sources/{id}", api.Source(ms)).Methods("PUT", "DELETE")
	r.HandleFunc("/roles", api.CreateRole(us)).Methods("POST")
	r.HandleFunc("/roles/{id}", api.UpdateRole(us)).Methods("PUT")
	r.HandleFunc("/roles/{id}", api.DeleteRole(us)).Methods("DELETE")
	r.HandleFunc("/steps", api.Steps(ls)).Methods("POST")
	r.HandleFunc("/steps/{id}", api.Step(ls)).Methods("PUT", "DELETE")
	r.HandleFunc("/products", api.Products(ms)).Methods("POST")
	r.HandleFunc("/products/{id}", api.Product(ms)).Methods("PUT", "DELETE")
	r.HandleFunc("/manufacturers", api.Manufacturers(ms)).Methods("POST")
	r.HandleFunc("/manufacturers/{id}", api.Manufacturer(ms)).Methods("PUT", "DELETE")
	// r.HandleFunc("/tags", api.TagsHandler).Methods("POST")
	// r.HandleFunc("/tags/{id}", api.TagHandler).Methods("PUT", "DELETE")
	r.HandleFunc("/tasks/{id}", api.Task(ts)).Methods("DELETE")
	// r.HandleFunc("/tasktypes", api.TaskTypesHandler).Methods("POST")
	// r.HandleFunc("/tasktypes/{id}", api.TaskTypeHandler).Methods("PUT", "DELETE")
	r.HandleFunc("/events/transfers", events_api.ListHandler).Methods("GET")
	r.HandleFunc("/files/{id}", api.File(fs)).Methods("DELETE")
	// Finance part
	r.HandleFunc("/wallets", api.Wallets(fin)).Methods("POST")
	r.HandleFunc("/wallets/{id}", api.Wallet(fin)).Methods("PUT", "DELETE")
	// TODO: changed to put - force front to use new route
	r.HandleFunc("/wallets/{id}/state", api.ChangeWalletState(fin)).Methods("GET")
	// r.HandleFunc("/wallets/{id}/open", api.OpenWalletHandler).Methods("GET")
	r.HandleFunc("/transfers/{id}", api.TransferHandler(fin)).Methods("DELETE")
	r.HandleFunc("/transfers/{id}/complete", api.CompleteTransferHandler(fin)).Methods("GET")
	return r
}

func UserRoutes(
	r *mux.Router, us users.Service, cs contacts.Service,
	ls leads.Service, ms misc.Service, ts tasks.Service,
	fs files.Service, fin finances.Service,
) *mux.Router {
	r.HandleFunc("/contacts", api.Contacts(cs)).Methods("GET", "POST")
	r.HandleFunc("/contacts/{id}", api.Contact(cs)).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/leads", api.Leads(ls)).Methods("GET", "POST")
	r.HandleFunc("/leads/{id}", api.Lead(ls)).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/users", api.Users(us)).Methods("GET")
	r.HandleFunc("/users/{id}", api.User(us)).Methods("GET")
	r.HandleFunc("/sources", api.Sources(ms)).Methods("GET")
	r.HandleFunc("/sources/{id}", api.Source(ms)).Methods("GET")
	r.HandleFunc("/roles", api.Roles(us)).Methods("GET")
	r.HandleFunc("/roles/{id}", api.Role(us)).Methods("GET")
	r.HandleFunc("/steps", api.Steps(ls)).Methods("GET")
	r.HandleFunc("/steps/{id}", api.Step(ls)).Methods("GET")
	r.HandleFunc("/products", api.Products(ms)).Methods("GET")
	r.HandleFunc("/products/{id}", api.Product(ms)).Methods("GET")
	r.HandleFunc("/manufacturers", api.Manufacturers(ms)).Methods("GET")
	r.HandleFunc("/manufacturers/{id}", api.Manufacturer(ms)).Methods("GET")
	// r.HandleFunc("/tags", api.TagsHandler).Methods("GET")
	// r.HandleFunc("/tags/{id}", api.TagHandler).Methods("GET")
	r.HandleFunc("/tasks", api.Tasks(ts)).Methods("GET", "POST")
	r.HandleFunc("/tasks/{id}", api.Task(ts)).Methods("GET", "PUT")
	// r.HandleFunc("/tasktypes", api.TaskTypesHandler).Methods("GET")
	// r.HandleFunc("/tasktypes/{id}", api.TaskTypeHandler).Methods("GET")
	r.HandleFunc("/files", api.Files(fs)).Methods("POST", "GET")
	r.HandleFunc("/files/{id}", api.File(fs)).Methods("GET")

	r.HandleFunc("/wallets", api.Wallets(fin)).Methods("GET")
	r.HandleFunc("/transfers", api.TransfersHandler(fin)).Methods("GET", "POST")
	r.HandleFunc("/transfers/{id}", api.TransferHandler(fin)).Methods("PUT")
	r.HandleFunc("/categories", api.CategoriesHandler(fin)).Methods("GET")
	r.HandleFunc("/analytics/categories", api.CategoriesSumHandler(fin)).Methods("GET")
	return r
}
