package routes

import (
	"github.com/gorilla/mux"
	"github.com/hromov/jevelina/api/events_api"
	"github.com/hromov/jevelina/auth"
)

func EventsRoutes(r *mux.Router) *mux.Router {
	r.HandleFunc("/events/transfers", auth.AdminCheck(events_api.ListHandler)).Methods("GET")
	return r
}
