package routes

import (
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/api"
	"github.com/hromov/jevelina/auth"
	"github.com/hromov/jevelina/domain/users"
)

func Base(us users.Service) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/usercheck", auth.UserCheckHandler).Methods("GET")
	r.HandleFunc("/orders", api.OrderHandler).Methods("POST")
	// TODO: uncoment for prod
	// r.Use(auth.UserCheck)
	r = UserRoutes(r, us)
	// TODO: uncoment for prod
	// r.Use(auth.AdminCheck)
	r = AdminRoutes(r, us)

	return r
}
