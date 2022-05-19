package routes

import (
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/api"
	"github.com/hromov/jevelina/auth"
)

func Base() *mux.Router {
	r := mux.NewRouter()
	r = UserRoutes(r)
	r = AdminRoutes(r)
	r = FinRoutes(r)
	r = FilesRoutes(r)
	r = EventsRoutes(r)
	r.HandleFunc("/usercheck", auth.UserCheckHandler).Methods("GET")
	r.HandleFunc("/orders", api.OrderHandler).Methods("POST")
	return r
}
