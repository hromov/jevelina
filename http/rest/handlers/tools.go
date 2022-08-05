package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getRouteID(r *http.Request) (uint64, error) {
	vars := mux.Vars(r)
	return strconv.ParseUint(vars["id"], 10, 32)
}
